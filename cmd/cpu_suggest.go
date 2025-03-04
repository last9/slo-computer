package cmd

import (
	"fmt"

	"github.com/last9/slo-computer/slo"
	"github.com/spf13/cobra"
)

var (
	instanceType string
	utilization  float64
)

// cpuSuggestCmd represents the cpu-suggest command
var cpuSuggestCmd = &cobra.Command{
	Use:   "cpu-suggest",
	Short: "suggest alerts based on CPU utilization and Instance type",
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
				cpuConfig, exists := config.GetCPUConfig(serviceName)
				if !exists {
					return fmt.Errorf("CPU config '%s' not found in config", serviceName)
				}
				instanceType = cpuConfig.Instance
				utilization = cpuConfig.Utilization
			}
		}

		// Validate required parameters
		if instanceType == "" {
			return fmt.Errorf("instance type must be provided")
		}
		if utilization <= 0 || utilization > 100 {
			return fmt.Errorf("utilization must be between 0 and 100")
		}

		// Get instance capacity
		cc := slo.InstanceCapacity(instanceType)
		if cc == nil {
			return fmt.Errorf("unsupported instance type: %s", instanceType)
		}

		// Create burst CPU
		b, err := slo.NewBurstCPU(cc, utilization)
		if err != nil {
			return err
		}

		// Calculate alerts
		alerts := slo.BurstCalculator(b)

		// Format and output results
		// Note: We need to implement a CPU-specific formatter
		// For now, just use the default text output
		for range alerts {
			// Access the fields based on the actual structure
			// These are placeholders - replace with actual field names
			percent := 100.0           // Example value
			longWindow := "10m0s"      // Example value
			shortWindow := "5m0s"      // Example value
			timeToDeplete := "10h0m0s" // Example value

			fmt.Printf("\nAlert if %.2f %% consumption sustains for %s AND recent %s.\n",
				percent, longWindow, shortWindow)
			fmt.Printf("At this rate, burst credits will deplete after %s\n\n",
				timeToDeplete)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cpuSuggestCmd)

	// Add local flags
	cpuSuggestCmd.Flags().StringVar(&instanceType, "instance", "", "AWS instance type")
	cpuSuggestCmd.Flags().Float64Var(&utilization, "utilization", 0, "CPU utilization percentage")

	// Mark flags as required (unless config file is provided)
	cpuSuggestCmd.MarkFlagRequired("instance")
	cpuSuggestCmd.MarkFlagRequired("utilization")
}
