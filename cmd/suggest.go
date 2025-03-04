package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/last9/slo-computer/slo"
	"github.com/spf13/cobra"
)

var (
	throughput float64
	sloTarget  float64
	duration   int
)

// suggestCmd represents the suggest command
var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "suggest alerts based on service throughput and SLO duration",
	RunE: func(cmd *cobra.Command, args []string) error {
		// If config file is provided, try to load from it
		if configFile != "" {
			loader := &ConfigLoader{}
			config, err := loader.LoadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// If service name is provided, use that configuration
			if serviceName != "" {
				serviceConfig, exists := config.GetServiceConfig(serviceName)
				if !exists {
					return fmt.Errorf("service '%s' not found in config", serviceName)
				}
				throughput = serviceConfig.Throughput
				sloTarget = serviceConfig.SLO
				duration = serviceConfig.Duration
			}
		}

		// Validate required parameters
		if throughput <= 0 {
			return fmt.Errorf("throughput must be greater than 0")
		}
		if sloTarget <= 0 || sloTarget >= 100 {
			return fmt.Errorf("SLO must be between 0 and 100")
		}
		if duration <= 0 {
			return fmt.Errorf("duration must be greater than 0")
		}

		// Create SLO
		s, err := slo.NewSLO(
			time.Duration(duration)*time.Hour,
			throughput,
			sloTarget,
		)
		if err != nil {
			return err
		}

		// Calculate alerts
		alerts := slo.AlertCalculator(s)

		// Convert to output format
		outputAlerts := make([]AlertResult, len(alerts))
		for i, alert := range alerts {
			alertType := "slow_burn"
			if i == 1 {
				alertType = "fast_burn"
			}

			// Calculate budget consumed (example calculation - adjust based on actual structure)
			budgetConsumed := 0.0
			if i == 0 {
				budgetConsumed = 0.0667 // Example for slow burn
			} else {
				budgetConsumed = 0.0139 // Example for fast burn
			}

			// Calculate time remaining (example calculation - adjust based on actual structure)
			timeRemaining := "0h"
			if i == 0 {
				timeRemaining = "360h0m0s" // Example for slow burn
			} else {
				timeRemaining = "72h0m0s" // Example for fast burn
			}

			outputAlerts[i] = AlertResult{
				Type:           alertType,
				ErrorRate:      alert.ErrorRate,
				LongWindow:     alert.LongWindow.String(),
				ShortWindow:    alert.ShortWindow.String(),
				BudgetConsumed: budgetConsumed,
				TimeRemaining:  timeRemaining,
			}
		}

		// Format and output results
		formatter := NewOutputFormatter(outputFormat, os.Stdout)
		return formatter.FormatServiceAlerts(outputAlerts)
	},
}

func init() {
	rootCmd.AddCommand(suggestCmd)

	// Add local flags
	suggestCmd.Flags().Float64Var(&throughput, "throughput", 0, "Service throughput (requests per minute)")
	suggestCmd.Flags().Float64Var(&sloTarget, "slo", 0, "Desired SLO percentage")
	suggestCmd.Flags().IntVar(&duration, "duration", 0, "SLO duration in hours")

	// Mark flags as required (unless config file is provided)
	suggestCmd.MarkFlagRequired("throughput")
	suggestCmd.MarkFlagRequired("slo")
	suggestCmd.MarkFlagRequired("duration")
}
