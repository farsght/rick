package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config represents the structure of your configuration
type Config struct {
	OpenAIKey     string `yaml:"openai_key"`
	AnthropicKey  string `yaml:"anthropic_key"`
	PerplexityKey string `yaml:"perplexity_key"`
	ActiveModel   string `yaml:"active_model"`
	Models        struct {
		OpenAI     string `yaml:"openai"`
		Anthropic  string `yaml:"anthropic"`
		Perplexity string `yaml:"perplexity"`
	} `yaml:"models"`
	Version  string `yaml:"version"`
	BuildNum string `yaml:"build_num"`
}

// Default configuration
var DefaultConfig = Config{
	OpenAIKey:     "",
	AnthropicKey:  "",
	PerplexityKey: "",
	ActiveModel:   "openai",
	Models: struct {
		OpenAI     string `yaml:"openai"`
		Anthropic  string `yaml:"anthropic"`
		Perplexity string `yaml:"perplexity"`
	}{
		OpenAI:     "gpt-4",
		Anthropic:  "claude-2",
		Perplexity: "mixtral-8x7b-instruct",
	},
	Version:  "1.0.0",
	BuildNum: "0",
}

// Helper function to mask API key
func MaskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

// GetConfigFilePath returns the full path to the config file
func GetConfigFilePath() string {
	_, configFile := GetConfigPath()
	return configFile
}

// GetConfigPath returns the path to the config directory and file
func GetConfigPath() (string, string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return "", ""
	}

	configDir := filepath.Join(homeDir, ".rick")
	configFile := filepath.Join(configDir, "rick.config.yaml")
	return configDir, configFile
}

func InitConfig() error {
	configDir, configFile := GetConfigPath()

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create and write default config
		data, err := yaml.Marshal(DefaultConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %v", err)
		}

		if err := os.WriteFile(configFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write config file: %v", err)
		}

		fmt.Printf("Created default configuration at %s\n", configFile)
	}

	return nil
}

// LoadConfig reads the configuration file
func LoadConfig() (*Config, error) {
	_, configFile := GetConfigPath()

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// SaveConfig writes the configuration to file
func SaveConfig(config *Config) error {
	_, configFile := GetConfigPath()

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
