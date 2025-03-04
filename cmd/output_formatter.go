package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// OutputFormat defines the format for command output
type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJSON OutputFormat = "json"
	OutputFormatYAML OutputFormat = "yaml"
)

// OutputFormatter handles formatting command results
type OutputFormatter struct {
	format OutputFormat
	writer io.Writer
}

// NewOutputFormatter creates a new formatter with the specified format
func NewOutputFormatter(format string, writer io.Writer) *OutputFormatter {
	outputFormat := OutputFormatText
	switch format {
	case "json":
		outputFormat = OutputFormatJSON
	case "yaml":
		outputFormat = OutputFormatYAML
	}

	return &OutputFormatter{
		format: outputFormat,
		writer: writer,
	}
}

// FormatServiceAlerts formats service SLO alerts in the configured format
func (f *OutputFormatter) FormatServiceAlerts(alerts []AlertResult) error {
	switch f.format {
	case OutputFormatJSON:
		return json.NewEncoder(f.writer).Encode(alerts)
	case OutputFormatYAML:
		return yaml.NewEncoder(f.writer).Encode(alerts)
	default:
		return f.formatServiceAlertsText(alerts)
	}
}

// formatServiceAlertsText formats alerts as human-readable text
func (f *OutputFormatter) formatServiceAlertsText(alerts []AlertResult) error {
	for _, alert := range alerts {
		fmt.Fprintf(f.writer, "\nAlert if error_rate > %.6f for last [%s] and also last [%s]\n",
			alert.ErrorRate, alert.LongWindow, alert.ShortWindow)
		fmt.Fprintf(f.writer, "This alert will trigger once %.2f%% of error budget is consumed,\n",
			alert.BudgetConsumed*100)
		fmt.Fprintf(f.writer, "and leaves %s before the SLO is defeated.\n\n",
			alert.TimeRemaining)
	}
	return nil
}

// AlertResult represents a structured alert recommendation
type AlertResult struct {
	Type           string  `json:"type" yaml:"type"`                       // "slow_burn" or "fast_burn"
	ErrorRate      float64 `json:"error_rate" yaml:"error_rate"`           // Error rate threshold
	LongWindow     string  `json:"long_window" yaml:"long_window"`         // Longer time window
	ShortWindow    string  `json:"short_window" yaml:"short_window"`       // Shorter time window
	BudgetConsumed float64 `json:"budget_consumed" yaml:"budget_consumed"` // Percentage of budget consumed
	TimeRemaining  string  `json:"time_remaining" yaml:"time_remaining"`   // Time until SLO breach
}
