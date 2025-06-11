package cmd

import (
	"encoding/json"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"os"
)

func InitLogger() error {
	return logger.InitLogger(&logger.LogCfg{
		IsStdout: true,
		File:     false,
	})
}

func readCfg(cfgFile string) (node.Config, error) {
	cfgBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		return node.Config{}, err
	}
	var runCfg node.RunConfig
	err = json.Unmarshal(cfgBytes, &runCfg)
	if err != nil {
		return node.Config{}, err
	}
	config, err := node.NewConfig(runCfg)
	if err != nil {
		return node.Config{}, err
	}
	return config, nil
}
