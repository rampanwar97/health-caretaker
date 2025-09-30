package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"health-monitoring/internal/models"
)

// Config represents the application configuration
type Config struct {
	Endpoints []EndpointConfig `json:"endpoints"`
	Server    ServerConfig     `json:"server"`
	Metrics   MetricsConfig    `json:"metrics"`
}

// EndpointConfig represents a single endpoint configuration
type EndpointConfig struct {
	Name      string            `json:"name"`
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Interval  int               `json:"interval"`
	Timeout   int               `json:"timeout"`
	Labels    map[string]string `json:"labels,omitempty"`     // Additional labels for metrics
	ProbeType string            `json:"probe_type,omitempty"` // e.g., "livez", "readyz", "healthz"
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port string `json:"port"`
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
	Port    string `json:"port"`
}

// LoadConfig loads configuration from a JSON file with environment variable overrides
func LoadConfig(filename string) (*Config, error) {
	var config *Config

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create default config if file doesn't exist
		config = getDefaultConfig()

		// Save default config
		if err := SaveConfig(filename, config); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
	} else {
		// Read existing config file
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}

		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %v", err)
		}
	}

	// Apply environment variable overrides
	applyEnvOverrides(config)

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %v", err)
	}

	return config, nil
}

// getDefaultConfig returns a default configuration
func getDefaultConfig() *Config {
	return &Config{
		Endpoints: []EndpointConfig{
			{
				Name:     "Google",
				URL:      "https://www.google.com",
				Method:   "GET",
				Interval: 30,
				Timeout:  10,
				Labels: map[string]string{
					"service":     "search",
					"environment": "production",
					"team":        "platform",
					"criticality": "high",
				},
			},
			{
				Name:     "GitHub",
				URL:      "https://api.github.com",
				Method:   "GET",
				Interval: 60,
				Timeout:  15,
				Labels: map[string]string{
					"service":     "api",
					"environment": "production",
					"team":        "development",
					"criticality": "medium",
				},
			},
		},
		Server: ServerConfig{
			Port: "8080",
		},
		Metrics: MetricsConfig{
			Enabled: true,
			Path:    "/metrics",
			Port:    "9091",
		},
	}
}

// applyEnvOverrides applies environment variable overrides to the configuration
func applyEnvOverrides(config *Config) {
	// Server configuration overrides
	if port := os.Getenv("WEB_PORT"); port != "" {
		config.Server.Port = port
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}

	// Metrics configuration overrides
	if enabled := os.Getenv("METRICS_ENABLED"); enabled != "" {
		config.Metrics.Enabled = enabled == "true"
	}

	if port := os.Getenv("METRICS_PORT"); port != "" {
		config.Metrics.Port = port
	}

	if path := os.Getenv("METRICS_PATH"); path != "" {
		config.Metrics.Path = path
	}
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(filename string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Metrics.Enabled {
		if c.Metrics.Port == "" {
			return fmt.Errorf("metrics port is required when metrics are enabled")
		}
		if c.Metrics.Path == "" {
			return fmt.Errorf("metrics path is required when metrics are enabled")
		}
	}

	for i, endpoint := range c.Endpoints {
		if err := endpoint.Validate(); err != nil {
			return fmt.Errorf("endpoint %d validation failed: %v", i, err)
		}
	}

	return nil
}

// Validate validates an endpoint configuration
func (ec *EndpointConfig) Validate() error {
	if ec.Name == "" {
		return fmt.Errorf("name is required")
	}

	if ec.URL == "" {
		return fmt.Errorf("URL is required")
	}

	if !strings.HasPrefix(ec.URL, "http://") && !strings.HasPrefix(ec.URL, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}

	if ec.Method == "" {
		ec.Method = "GET"
	}

	if ec.Interval <= 0 {
		ec.Interval = 30
	}

	if ec.Timeout <= 0 {
		ec.Timeout = 10
	}

	return nil
}

// ToEndpoint converts EndpointConfig to models.Endpoint
func (ec *EndpointConfig) ToEndpoint() *models.Endpoint {
	return &models.Endpoint{
		Name:      ec.Name,
		URL:       ec.URL,
		Method:    ec.Method,
		Interval:  ec.Interval,
		Timeout:   ec.Timeout,
		Status:    "checking",
		Labels:    ec.Labels,
		ProbeType: ec.ProbeType,
	}
}
