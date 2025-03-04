package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile   string
	outputFormat string
	serviceName  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slo-computer",
	Short: "Last9 SLO toolkit",
	Long:  `A toolkit for calculating SLO-based alert thresholds for services and AWS burstable instances.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to configuration file")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "output", "text", "Output format (text, json, yaml)")
	rootCmd.PersistentFlags().StringVar(&serviceName, "service", "", "Service name to use from config file")
}
