package cmd

import (
	dcommon "github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"testing"
)

func TestBeaconNode_Check(t *testing.T) {
	config := Config{}
	err := dcommon.ReadObj("/Users/red/lworkspace/lightec/audit/daemon/cmd/monitor/config.json", &config)
	if err != nil {
		panic(err)
	}
	err = logger.InitLogger(&logger.LogCfg{
		IsStdout:       true,
		DiscordHookUrl: config.DiscordUrl,
	})
	if err != nil {
		panic(err)
	}
	node, err := NewNode(config)
	if err != nil {
		panic(err)
	}
	err = node.Run()
	if err != nil {
		panic(err)
	}
}
