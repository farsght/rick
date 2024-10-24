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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Rick",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration to get version info
		config, err := LoadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
		fmt.Printf("Rick Version: %s (Build %s)\n", config.Version, config.BuildNum)
	},
}

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "LLM chatbot",
	Long:  `Rick is a silly goose ai chatbot.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		config, err := LoadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		// Check for API key and prompt if not found
		if config.APIKey == "" {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("\n🔑 No API key found in configuration.\n")
			fmt.Print("Please enter your OpenAI API key: ")

			apiKey, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading API key: %v", err)
			}

			// Clean the input
			apiKey = strings.TrimSpace(apiKey)

			if apiKey == "" {
				log.Fatal("API key cannot be empty")
			}

			// Update config with new API key
			config.APIKey = apiKey

			// Save the updated configuration
			if err := SaveConfig(config); err != nil {
				log.Fatalf("Error saving configuration: %v", err)
			}

			fmt.Printf("\n✅ API key saved to %s\n\n", GetConfigFilePath())
		}

		// Rest of your chat command implementation...
		reader := bufio.NewReader(os.Stdin)

		// Set up interrupt handler
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			fmt.Println("\nInterrupt signal received. Exiting...")
			os.Exit(0)
		}()

		// Create OpenAI client
		llm, err := openai.New(
			openai.WithToken(config.APIKey),
			openai.WithModel(config.Model),
		)
		if err != nil {
			log.Fatalf("Error creating OpenAI client: %v", err)
		}

		ctx := context.Background()

		// Initial LLM prompt phase
		initalPromptColor.Print("Optional Rick Context: ")
		initialPrompt, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}

		systemPrompt := "You are an AI assistant named Rick. Always refer to yourself as Rick when appropriate."
		content := []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		}

		if strings.TrimSpace(initialPrompt) != "" {
			content = append(content, llms.TextParts(llms.ChatMessageTypeSystem, initialPrompt))
		}

		fmt.Println("Initial prompt received. Entering chat mode...")

		// Main chat loop
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
				response := ""
				content = append(content, llms.TextParts(llms.ChatMessageTypeHuman, input))
				llmOutputColor.Print("[RICK] ")

				_, err = llm.GenerateContent(ctx, content,
					llms.WithMaxTokens(1024),
					llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
						llmOutputColor.Print(string(chunk))
						response += string(chunk)
						return nil
					}),
				)

				if err != nil {
					log.Printf("Error generating response: %v", err)
					continue
				}

				content = append(content, llms.TextParts(llms.ChatMessageTypeSystem, response))
				fmt.Println()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(versionCmd)
}
