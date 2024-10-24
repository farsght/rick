package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"rick/internal/config" // Update this import path to match your module name
	"rick/internal/utils"  // For cleanMarkdown and other utilities

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

var (
	userInputColor    = color.New(color.FgGreen, color.Bold)
	llmOutputColor    = color.New(color.FgMagenta, color.Bold)
	initalPromptColor = color.New(color.FgYellow, color.Bold)
)

// Helper function to prompt for API key
func promptForAPIKey(provider string, cfg *config.Config) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n🔑 No %s API key found in configuration.\n", provider)
	fmt.Printf("Please enter your %s API key: ", provider)

	apiKey, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading API key: %v", err)
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		log.Fatal("API key cannot be empty")
	}

	// Save the API key to the correct field
	switch provider {
	case "OpenAI":
		cfg.OpenAIKey = apiKey
	case "Anthropic":
		cfg.AnthropicKey = apiKey
	case "Perplexity":
		cfg.PerplexityKey = apiKey
	}

	if err := config.SaveConfig(cfg); err != nil {
		log.Fatalf("Error saving configuration: %v", err)
	}

	fmt.Printf("\n✅ %s API key saved to %s\n\n", provider, config.GetConfigFilePath())
	return apiKey
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "LLM chatbot",
	Long:  `Rick is a silly goose ai chatbot.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		// Set default provider if none is set
		if cfg.ActiveModel == "" {
			cfg.ActiveModel = "openai"
			if err := config.SaveConfig(cfg); err != nil {
				log.Fatalf("Error saving config: %v", err)
			}
			log.Printf("No active model provider set, using default: %s", cfg.ActiveModel)
		}

		// Initialize LLM client based on provider
		var llm *openai.LLM

		switch cfg.ActiveModel {
		case "openai":
			if cfg.OpenAIKey == "" {
				cfg.OpenAIKey = promptForAPIKey("OpenAI", cfg)
			}
			if cfg.Models.OpenAI == "" {
				cfg.Models.OpenAI = "gpt-4"
				log.Printf("No model specified for OpenAI, using default: %s", cfg.Models.OpenAI)
			}
			llm, err = openai.New(
				openai.WithToken(cfg.OpenAIKey),
				openai.WithModel(cfg.Models.OpenAI),
			)

		case "perplexity":
			if cfg.PerplexityKey == "" {
				cfg.PerplexityKey = promptForAPIKey("Perplexity", cfg)
			}
			if cfg.Models.Perplexity == "" {
				cfg.Models.Perplexity = "llama-3.1-sonar-large-128k-online"
				log.Printf("No model specified for Perplexity, using default: %s", cfg.Models.Perplexity)
			}
			llm, err = openai.New(
				openai.WithToken(cfg.PerplexityKey),
				openai.WithModel(cfg.Models.Perplexity),
				openai.WithBaseURL("https://api.perplexity.ai"),
			)

		case "anthropic":
			if cfg.AnthropicKey == "" {
				cfg.AnthropicKey = promptForAPIKey("Anthropic", cfg)
			}
			if cfg.Models.Anthropic == "" {
				cfg.Models.Anthropic = "claude-2"
				log.Printf("No model specified for Anthropic, using default: %s", cfg.Models.Anthropic)
			}
			log.Fatalf("Anthropic support not yet implemented")

		default:
			log.Fatalf("Unknown model provider: %s", cfg.ActiveModel)
		}

		if err != nil {
			log.Fatalf("Error creating LLM client: %v", err)
		}

		// Get current model name for display
		var currentModel string
		switch cfg.ActiveModel {
		case "openai":
			currentModel = cfg.Models.OpenAI
		case "perplexity":
			currentModel = cfg.Models.Perplexity
		case "anthropic":
			currentModel = cfg.Models.Anthropic
		}

		fmt.Printf("\nUsing %s model: %s\n", cfg.ActiveModel, currentModel)

		// Set up chat session
		reader := bufio.NewReader(os.Stdin)
		setupInterruptHandler()
		ctx := context.Background()

		// Initialize chat content
		content := initializeChatContent(cfg.ActiveModel, reader)

		// Main chat loop
		runChatLoop(ctx, llm, content, cfg.ActiveModel)
	},
}

func setupInterruptHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nInterrupt signal received. Exiting...")
		os.Exit(0)
	}()
}

func initializeChatContent(activeModel string, reader *bufio.Reader) []llms.MessageContent {
	systemPrompt := `You are an AI assistant named Rick. Always refer to yourself as Rick when appropriate. 
    IMPORTANT FORMATTING INSTRUCTIONS:
    - Use plain text only
    - Do not use any markdown syntax or special formatting
    - Do not use asterisks, hashtags, or any other markdown symbols
    - Do not use HTML tags
    - Just respond with natural, plain text`

	initalPromptColor.Print("Optional Rick Context: ")
	initialPrompt, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	var content []llms.MessageContent
	if activeModel == "perplexity" {
		content = initializePerplexityContent(systemPrompt, initialPrompt)
	} else {
		content = initializeStandardContent(systemPrompt, initialPrompt)
	}

	fmt.Println("Initial prompt received. Entering chat mode...")
	return content
}

func initializePerplexityContent(systemPrompt, initialPrompt string) []llms.MessageContent {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, "Hello Rick, introduce yourself."),
		llms.TextParts(llms.ChatMessageTypeAI, "Hi! I'm Rick, your AI assistant. How can I help you today?"),
	}

	if strings.TrimSpace(initialPrompt) != "" {
		content = append(content,
			llms.TextParts(llms.ChatMessageTypeHuman, initialPrompt),
			llms.TextParts(llms.ChatMessageTypeAI, "I understand. How can I help you?"),
		)
	}
	return content
}

func initializeStandardContent(systemPrompt, initialPrompt string) []llms.MessageContent {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
	}
	if strings.TrimSpace(initialPrompt) != "" {
		content = append(content, llms.TextParts(llms.ChatMessageTypeSystem, initialPrompt))
	}
	return content
}

func runChatLoop(ctx context.Context, llm *openai.LLM, content []llms.MessageContent, activeModel string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		userInputColor.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}

		input = strings.TrimSpace(input)

		switch input {
		case "quit", "exit":
			fmt.Println("Exiting...")
			os.Exit(0)
		case "":
			continue
		default:
			handleChatResponse(ctx, llm, &content, input, activeModel)
		}
	}
}

func handleChatResponse(ctx context.Context, llm *openai.LLM, content *[]llms.MessageContent, input, activeModel string) {
	response := ""
	*content = append(*content, llms.TextParts(llms.ChatMessageTypeHuman, input))
	llmOutputColor.Print("[RICK] ")

	_, err := llm.GenerateContent(ctx, *content,
		llms.WithMaxTokens(1024),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			cleanedChunk := utils.CleanMarkdown(string(chunk))
			llmOutputColor.Print(cleanedChunk)
			response += cleanedChunk
			return nil
		}),
	)

	if err != nil {
		log.Printf("Error generating response: %v", err)
		return
	}

	if activeModel == "perplexity" {
		if len(*content) > 6 {
			*content = (*content)[len(*content)-6:]
		}
		*content = append(*content, llms.TextParts(llms.ChatMessageTypeAI, response))
	} else {
		*content = append(*content, llms.TextParts(llms.ChatMessageTypeSystem, response))
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
