package cmd

import (
	"context"
	"fmt"
	dcommon "github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	dnode "github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Printf("get path error: %v \n", err)
			return
		}
		config := Config{}
		err = dcommon.ReadObj(path, &config)
		if err != nil {
			fmt.Printf("read config error: %v \n", err)
			return
		}
		err = logger.InitLogger(&logger.LogCfg{
			DiscordHookUrl: config.DiscordUrl,
			IsStdout:       true,
		})
		if err != nil {
			fmt.Printf("init logger error: %v \n", err)
			return
		}

		node, err := NewNode(config)
		if err != nil {
			fmt.Printf("new node error: %v \n", err)
			return
		}
		err = node.Run()
		if err != nil {
			fmt.Printf("run node error: %v \n", err)
			return
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("config", "./config.json", "rpc node urls config")
}

type Node struct {
	nodes []INode
	exit  chan os.Signal
}

func NewNode(cfg Config) (*Node, error) {
	var nodes []INode
	for _, url := range cfg.EthUrls {
		node, err := NewEthNode(url.Url)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	for _, url := range cfg.BtcUrls {
		node, err := NewBtcNode(url.Url, url.User, url.Pwd)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	for _, url := range cfg.BeaconUrls {
		node, err := NewBeaconNode(url.Url)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return &Node{
		nodes: nodes,
		exit:  make(chan os.Signal, 1),
	}, nil
}

func (n *Node) Run() error {
	for _, tn := range n.nodes {
		go dnode.DoTimerTask(fmt.Sprintf("monitor-%v", tn.Name()), tn.Time(), tn.Check, n.exit)
	}
	logger.Debug("monitor is running now ...")
	signal.Notify(n.exit, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		msg := <-n.exit
		switch msg {
		case syscall.SIGQUIT, syscall.SIGTERM:
			logger.Info("get shutdown signal , exit now ...")
			err := n.Close()
			if err != nil {
				logger.Error("%v", err)
			}
			return nil
		}
	}
}

func (n *Node) Close() error {
	return nil
}

type EthNode struct {
	curHeight uint64
	client    *ethereum.Client
	time      time.Duration
	url       string
}

func NewEthNode(url string) (*EthNode, error) {
	client, err := ethereum.NewClient(url, "", "", "", "")
	if err != nil {
		return nil, err
	}
	number, err := client.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	return &EthNode{
		curHeight: number,
		client:    client,
		time:      10 * time.Minute,
		url:       url,
	}, nil
}

func (e *EthNode) Check() error {
	number, err := e.client.BlockNumber(context.Background())
	if err != nil {
		logger.Error("get eth block number error:%v", err)
		return err
	}
	if e.curHeight != 0 {
		diff := number - e.curHeight
		if diff <= 35 { // 50 blocks per 10 minutes
			logger.Error("ethereum sync too slow, %v node maybe offline: diff %v prevHeight:%v curHeight:%v", e.url, diff, e.curHeight, number)
		}
	}
	e.curHeight = number
	return nil
}

func (e *EthNode) Time() time.Duration {
	return e.time
}

func (e *EthNode) Name() string {
	return e.url
}

type BtcNode struct {
	curHeight int64
	client    *bitcoin.Client
	time      time.Duration
	url       string
}

func NewBtcNode(url, user, pwd string) (*BtcNode, error) {
	client, err := bitcoin.NewClient(url, user, pwd)
	if err != nil {
		return nil, err
	}
	height, err := client.GetBlockCount()
	if err != nil {
		return nil, err
	}
	return &BtcNode{
		curHeight: int64(int(height)),
		client:    client,
		time:      30 * time.Minute,
		url:       url,
	}, nil
}

func (b *BtcNode) Check() error {
	height, err := b.client.GetBlockCount()
	if err != nil {
		return err
	}
	if b.curHeight != 0 {
		diff := height - b.curHeight
		if diff <= 1 { // 3 blocks per 30 minutes
			logger.Error("bitcoin sync too slow,%v node maybe offline: diff %v prevHeight:%v curHeight:%v", b.url, diff, b.curHeight, height)
		}
	}
	b.curHeight = height
	return nil
}

func (b *BtcNode) Time() time.Duration {
	return b.time
}

func (b *BtcNode) Name() string {
	return b.url
}

type BeaconNode struct {
	curHeight uint64
	client    *beacon.Client
	time      time.Duration
	url       string
}

func NewBeaconNode(url string) (*BeaconNode, error) {
	client, err := beacon.NewClient(url)
	if err != nil {
		return nil, err
	}
	height, err := client.GetLatestFinalizedSlot()
	if err != nil {
		return nil, err
	}
	return &BeaconNode{
		curHeight: height,
		client:    client,
		time:      10 * time.Minute,
		url:       url,
	}, nil
}

func (b *BeaconNode) Check() error {
	height, err := b.client.GetLatestFinalizedSlot()
	if err != nil {
		return err
	}
	if b.curHeight != 0 {
		diff := height - b.curHeight
		if diff <= 35 { // 50 blocks per 10 minutes
			logger.Error("beacon sync too slow, %v node maybe offline: diff %v prevHeight:%v curHeight:%v", b.url, diff, b.curHeight, height)
		}
	}
	b.curHeight = height
	return nil
}

func (b *BeaconNode) Time() time.Duration {
	return b.time
}

func (b *BeaconNode) Name() string {
	return b.url
}

type Config struct {
	EthUrls    []Url  `json:"EthUrls"`
	BtcUrls    []Url  `json:"BtcUrls"`
	BeaconUrls []Url  `json:"BeaconUrls"`
	DiscordUrl string `json:"discordUrl"`
}
type Url struct {
	Url  string `json:"url"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

type INode interface {
	Check() error
	Name() string
	Time() time.Duration
}
