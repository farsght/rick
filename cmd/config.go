package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Config represents the structure of your configuration
type Config struct {
	APIKey   string `yaml:"api_key"`
	Model    string `yaml:"model"`
	Version  string `yaml:"version"`
	BuildNum string `yaml:"build_num"`
}

// Default configuration
var defaultConfig = Config{
	APIKey:   "",
	Model:    "gpt-4o",
	Version:  "1.0.0",
	BuildNum: "0",
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Rick configuration",
	Long: `Manage Rick configuration settings.
Example: rick config set api_key YOUR_API_KEY`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		if err := cmd.Help(); err != nil {
			log.Fatalf("Error showing help: %v", err)
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		config, err := LoadConfig()
		if err != nil {
			// Try to initialize config if it doesn't exist
			if err := InitConfig(); err != nil {
				log.Fatalf("Error initializing config: %v", err)
			}
			config = &defaultConfig
		}

		switch key {
		case "api_key":
			if value == "" {
				log.Fatal("API key cannot be empty")
			}
			config.APIKey = value
		case "model":
			if value == "" {
				log.Fatal("Model cannot be empty")
			}
			config.Model = value
		default:
			fmt.Printf("Available configuration keys:\n")
			fmt.Printf("  - api_key\n")
			fmt.Printf("  - model\n")
			log.Fatalf("Unknown configuration key: %s", key)
		}

		if err := SaveConfig(config); err != nil {
			log.Fatalf("Error saving config: %v", err)
		}

		fmt.Printf("✅ Set %s in configuration\n", key)

		// Show the updated configuration
		fmt.Println("\nUpdated Configuration:")
		if key == "api_key" {
			fmt.Printf("API Key: %s\n", maskAPIKey(config.APIKey))
		} else {
			fmt.Printf("%s: %s\n", key, value)
		}
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadConfig()
		if err != nil {
			// Try to initialize config if it doesn't exist
			if err := InitConfig(); err != nil {
				log.Fatalf("Error initializing config: %v", err)
			}
			config = &defaultConfig
		}

		fmt.Println("\nCurrent Configuration:")
		fmt.Printf("API Key: %s\n", maskAPIKey(config.APIKey))
		fmt.Printf("Model: %s\n", config.Model)
		fmt.Printf("Version: %s\n", config.Version)
		fmt.Printf("Build: %s\n", config.BuildNum)
		fmt.Printf("\nConfig file location: %s\n", GetConfigFilePath())
	},
}

// Helper function to mask API key
func maskAPIKey(key string) string {
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
		data, err := yaml.Marshal(defaultConfig)
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

func init() {
	// Add config commands to root command
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
