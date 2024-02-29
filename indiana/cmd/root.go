package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "indiana",
	Short: "CLI for Indiana container orchestration tool",
	Long: `indiana is the CLI for Indiana container orchestrator. 
	Indiana can create a highly available container infrastructure, based on Docker container platform.
	CLI functionalities to be developed:
		- separate cluster element changes (containers, networks, volumes)
		- log access over CLI. For now they are just stored in files`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize cluster by given yaml config file",
	Run: func(cmd *cobra.Command, args []string) {
		confFile, _ := cmd.Flags().GetString("config")
		fmt.Println("Local flag value:", confFile)
		f, err := os.Open(confFile)
		if err != nil {
			slog.Error("Could not open config file", "name", confFile)
		}
		defer f.Close()
		URL := "http://localhost:1986/clusterState"
		req, err := http.NewRequest(http.MethodPost, URL, f)
		if err != nil {
			slog.Error("Could not create POST request", "URL", URL, "err", err.Error())
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("Could not send POST request", "URL", URL, "err", err.Error())
			return
		}

		if resp.StatusCode == http.StatusOK {
			slog.Info("Cluster Change Outcome logs send successfully")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	initCmd.Flags().String("config", "/home/indiana/orchestrator/src/config/clusterState.yaml", "Pass yaml configuration file name")
	rootCmd.AddCommand(initCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".indiana")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
