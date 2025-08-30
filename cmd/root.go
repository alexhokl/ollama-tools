// Package cmd implements the CLI commands.
package cmd

import (
	"github.com/alexhokl/helper/cli"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "ollama-tools",
	Short:        "Generate function calls using Ollama",
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// ignore the error as error message should be printed by cobra
		return
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ollama-tools.yml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	cli.ConfigureViper(cfgFile, "ollama-tools", false, "")
}
