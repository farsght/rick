package cmd

import (
	"fmt"
	"log"
	"rick/internal/config" // Make sure this import path matches your module name

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Rick configuration",
	Long: `Manage Rick configuration settings.
Example: rick config set openai_key YOUR_API_KEY`,
	Run: func(cmd *cobra.Command, args []string) {
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

		cfg, err := config.LoadConfig()
		if err != nil {
			if err := config.InitConfig(); err != nil {
				log.Fatalf("Error initializing config: %v", err)
			}
			cfg = &config.DefaultConfig
		}

		switch key {
		case "openai_key":
			cfg.OpenAIKey = value
		case "anthropic_key":
			cfg.AnthropicKey = value
		case "perplexity_key":
			cfg.PerplexityKey = value
		case "model":
			cfg.ActiveModel = value
		default:
			fmt.Printf("Available configuration keys:\n")
			fmt.Printf("  - openai_key\n")
			fmt.Printf("  - anthropic_key\n")
			fmt.Printf("  - perplexity_key\n")
			fmt.Printf("  - model\n")
			log.Fatalf("Unknown configuration key: %s", key)
		}

		if err := cfg.SaveConfig(); err != nil {
			log.Fatalf("Error saving config: %v", err)
		}

		fmt.Printf("✅ Set %s in configuration\n", key)
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			if err := config.InitConfig(); err != nil {
				log.Fatalf("Error initializing config: %v", err)
			}
			cfg = &config.DefaultConfig
		}

		fmt.Println("\nCurrent Configuration:")
		fmt.Printf("OpenAI Key: %s\n", config.MaskAPIKey(cfg.OpenAIKey))
		fmt.Printf("Anthropic Key: %s\n", config.MaskAPIKey(cfg.AnthropicKey))
		fmt.Printf("Perplexity Key: %s\n", config.MaskAPIKey(cfg.PerplexityKey))
		fmt.Printf("Active Model: %s\n", cfg.ActiveModel)
		fmt.Printf("Version: %s\n", cfg.Version)
		fmt.Printf("Build: %s\n", cfg.BuildNum)
		fmt.Printf("\nConfig file location: %s\n", config.GetConfigFilePath())
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
