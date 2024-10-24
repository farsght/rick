package cmd

import (
	"fmt"
	"log"
	"rick/internal/config" // Update this import path to match your module name
	"strings"

	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage LLM models",
	Long: `Manage Large Language Model settings.
Examples:
  rick model list
  rick model set openai gpt-4
  rick model active openai`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatalf("Error showing help: %v", err)
		}
	},
}

var modelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available models",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		fmt.Println("\nAvailable Models:")
		fmt.Printf("OpenAI     [%s]: %s\n",
			getActiveIndicator(cfg, "openai"),
			cfg.Models.OpenAI)
		fmt.Printf("Anthropic  [%s]: %s\n",
			getActiveIndicator(cfg, "anthropic"),
			cfg.Models.Anthropic)
		fmt.Printf("Perplexity [%s]: %s\n",
			getActiveIndicator(cfg, "perplexity"),
			cfg.Models.Perplexity)
		fmt.Printf("\nActive Model: %s\n", cfg.ActiveModel)
	},
}

var modelSetCmd = &cobra.Command{
	Use:   "set [provider] [model]",
	Short: "Set model for a specific provider",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		provider := strings.ToLower(args[0])
		model := args[1]

		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		switch provider {
		case "openai":
			cfg.Models.OpenAI = model
		case "anthropic":
			cfg.Models.Anthropic = model
		case "perplexity":
			cfg.Models.Perplexity = model
		default:
			log.Fatalf("Unknown provider: %s. Available providers: openai, anthropic, perplexity", provider)
		}

		if err := config.SaveConfig(cfg); err != nil {
			log.Fatalf("Error saving config: %v", err)
		}

		fmt.Printf("✅ Set %s model to: %s\n", provider, model)
	},
}

var modelActiveCmd = &cobra.Command{
	Use:   "active [provider]",
	Short: "Set the active model provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider := strings.ToLower(args[0])

		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		switch provider {
		case "openai", "anthropic", "perplexity":
			cfg.ActiveModel = provider
		default:
			log.Fatalf("Unknown provider: %s. Available providers: openai, anthropic, perplexity", provider)
		}

		if err := config.SaveConfig(cfg); err != nil {
			log.Fatalf("Error saving config: %v", err)
		}

		fmt.Printf("✅ Active model set to: %s\n", provider)

		// Show current model for the active provider
		var currentModel string
		switch provider {
		case "openai":
			currentModel = cfg.Models.OpenAI
		case "anthropic":
			currentModel = cfg.Models.Anthropic
		case "perplexity":
			currentModel = cfg.Models.Perplexity
		}
		fmt.Printf("Current %s model: %s\n", provider, currentModel)
	},
}

// Helper function to show which model is active
func getActiveIndicator(cfg *config.Config, provider string) string {
	if cfg.ActiveModel == provider {
		return "*"
	}
	return " "
}

func init() {
	modelCmd.AddCommand(modelListCmd)
	modelCmd.AddCommand(modelSetCmd)
	modelCmd.AddCommand(modelActiveCmd)
	rootCmd.AddCommand(modelCmd)
}
