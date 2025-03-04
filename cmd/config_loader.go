package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// ConfigLoader handles loading configuration from files
type ConfigLoader struct{}

// ServiceConfig represents configuration for service SLO calculation
type ServiceConfig struct {
	Throughput float64 `json:"throughput" yaml:"throughput"`
	SLO        float64 `json:"slo" yaml:"slo"`
	Duration   int     `json:"duration" yaml:"duration"`
}

// CPUConfig represents configuration for CPU burst calculation
type CPUConfig struct {
	Instance    string  `json:"instance" yaml:"instance"`
	Utilization float64 `json:"utilization" yaml:"utilization"`
}

// Config represents the full configuration file structure
type Config struct {
	Services map[string]ServiceConfig `json:"services" yaml:"services"`
	CPUs     map[string]CPUConfig     `json:"cpus" yaml:"cpus"`
}

// LoadConfig loads configuration from a file
func (l *ConfigLoader) LoadConfig(path string) (*Config, error) {
	ext := filepath.Ext(path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	return config, nil
}

// GetServiceConfig retrieves a specific service configuration
func (c *Config) GetServiceConfig(name string) (ServiceConfig, bool) {
	if c.Services == nil {
		return ServiceConfig{}, false
	}

	config, exists := c.Services[name]
	return config, exists
}

// GetCPUConfig retrieves a specific CPU configuration
func (c *Config) GetCPUConfig(name string) (CPUConfig, bool) {
	if c.CPUs == nil {
		return CPUConfig{}, false
	}

	config, exists := c.CPUs[name]
	return config, exists
}
