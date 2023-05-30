package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	defaultConfigFile = "pow-fiat-shamir-config.yml"
)

var (
	homeDir string

	config *RuntimeConfig

	RootCmd = &cobra.Command{
		Use:   "pow-fiat-shamir",
		Short: "multi-sign needed for send eth transaction with multi sign",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	config = &RuntimeConfig{}
	handleInitError(configFile(config))
}

func handleInitError(err error) {
	if err != nil {
		fmt.Println("init error:", err)
		os.Exit(1)
	}
}

func configFile(config *RuntimeConfig) error {
	var home string
	if homeDir == "" {
		userHome, err := homedir.Dir()
		if err != nil {
			return err
		}
		home = filepath.Join(userHome, defaultConfigFile)
	} else {
		home = homeDir
	}

	bz, err := os.ReadFile(home)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bz, &config)
	if err != nil {
		return err
	}
	return nil
}
