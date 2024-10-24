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

	"rick/internal/config"
	"rick/internal/llm"
	"rick/internal/utils"

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
func promptForAPIKey(provider string, config *config.Config) string {
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
		config.OpenAIKey = apiKey
	case "Anthropic":
		config.AnthropicKey = apiKey
	case "Perplexity":
		config.PerplexityKey = apiKey
	}

	if err := config.SaveConfig(); err != nil {
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
		config, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		// Set default provider if none is set
		if config.ActiveModel == "" {
			config.ActiveModel = "openai"
			if err := config.SaveConfig(); err != nil { // Removed the argument
				log.Fatalf("Error saving config: %v", err)
			}
			log.Printf("No active model provider set, using default: %s", config.ActiveModel)
		}

		// Initialize LLM client based on provider
		var llm *openai.LLM

		switch config.ActiveModel {
		case "openai":
			if config.OpenAIKey == "" {
				log.Fatalf("No OpenAI API key found in configuration")
			}
			if config.Models.OpenAI == "" {
				config.Models.OpenAI = "gpt-4"
				log.Printf("No model specified for OpenAI, using default: %s", config.Models.OpenAI)
			}
			llm, err = openai.New(
				openai.WithToken(config.OpenAIKey),
				openai.WithModel(config.Models.OpenAI),
			)

		case "perplexity":
			if config.PerplexityKey == "" {
				log.Fatalf("No Perplexity API key found in configuration")
			}
			if config.Models.Perplexity == "" {
				config.Models.Perplexity = "llama-3.1-sonar-large-128k-online"
				log.Printf("No model specified for Perplexity, using default: %s", config.Models.Perplexity)
			}
			llm, err = openai.New(
				openai.WithToken(config.PerplexityKey),
				openai.WithModel(config.Models.Perplexity),
				openai.WithBaseURL("https://api.perplexity.ai"),
			)

		case "anthropic":
			if config.AnthropicKey == "" {
				log.Fatalf("No Anthropic API key found in configuration")
			}
			if config.Models.Anthropic == "" {
				config.Models.Anthropic = "claude-2"
				log.Printf("No model specified for Anthropic, using default: %s", config.Models.Anthropic)
			}
			// TODO: Implement Anthropic client configuration
			log.Fatalf("Anthropic support not yet implemented")

		default:
			log.Fatalf("Unknown model provider: %s", config.ActiveModel)
		}

		if err != nil {
			log.Fatalf("Error creating LLM client: %v", err)
		}

		// Get current model name for display
		var currentModel string
		switch config.ActiveModel {
		case "openai":
			currentModel = config.Models.OpenAI
		case "perplexity":
			currentModel = config.Models.Perplexity
		case "anthropic":
			currentModel = config.Models.Anthropic
		}

		fmt.Printf("\nUsing %s model: %s\n", config.ActiveModel, currentModel)

		// Set up chat session
		reader := bufio.NewReader(os.Stdin)
		setupInterruptHandler()
		ctx := context.Background()

		// Initialize chat content based on provider
		var content []llms.MessageContent
		if config.ActiveModel == "perplexity" {
			content = initializePerplexityContent(getSystemPrompt(), getInitialPrompt(reader))
		} else {
			content = initializeStandardContent(getSystemPrompt(), getInitialPrompt(reader))
		}

		// Main chat loop
		runChatLoop(ctx, llm, content, config.ActiveModel)
	},
}

func getSystemPrompt() string {
	return `You are an AI assistant named Rick. Always refer to yourself as Rick when appropriate. 
    IMPORTANT FORMATTING INSTRUCTIONS:
    - Use plain text only
    - Do not use any markdown syntax or special formatting
    - Do not use asterisks, hashtags, or any other markdown symbols
    - Do not use HTML tags
    - Just respond with natural, plain text`
}

func getInitialPrompt(reader *bufio.Reader) string {
	initalPromptColor.Print("Optional Rick Context: ")
	initialPrompt, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	return initialPrompt
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

func initializeStandardContent(systemPrompt, initialPrompt string) []llms.MessageContent {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
	}
	if strings.TrimSpace(initialPrompt) != "" {
		content = append(content, llms.TextParts(llms.ChatMessageTypeSystem, initialPrompt))
	}
	return content
}

func initializePerplexityContent(systemPrompt, initialPrompt string) []llms.MessageContent {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
	}
	if strings.TrimSpace(initialPrompt) != "" {
		content = append(content, llms.TextParts(llms.ChatMessageTypeSystem, initialPrompt))
	}
	return content
}

func runChatLoop(ctx context.Context, client *openai.LLM, content []llms.MessageContent, activeModel string) {
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
			if activeModel == "perplexity" {
				llm.HandlePerplexityResponse(&content, input)
			} else {
				handleStandardResponse(ctx, client, &content, input)
			}
		}
	}
}

func handleStandardResponse(ctx context.Context, client *openai.LLM, content *[]llms.MessageContent, input string) {
	response := ""
	*content = append(*content, llms.TextParts(llms.ChatMessageTypeHuman, input))
	llmOutputColor.Print("[RICK] ")

	_, err := client.GenerateContent(ctx, *content,
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

	*content = append(*content, llms.TextParts(llms.ChatMessageTypeSystem, response))
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
