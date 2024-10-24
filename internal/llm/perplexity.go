package llm

import (
	"rick/internal/config"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

const (
	DefaultPerplexityModel = "llama-3.1-sonar-large-128k-online"
	BaseURL                = "https://api.perplexity.ai"
)

// PerplexityProvider handles Perplexity-specific LLM operations
type PerplexityProvider struct {
	client *openai.LLM
}

// NewPerplexityClient creates a new Perplexity client
func NewPerplexityClient(cfg *config.Config) (*openai.LLM, error) {
	return openai.New(
		openai.WithToken(cfg.PerplexityKey),
		openai.WithModel(cfg.Models.Perplexity),
		openai.WithBaseURL(BaseURL),
	)
}

// InitializePerplexityContent creates the initial message content for Perplexity
func InitializePerplexityContent(systemPrompt, initialPrompt string) []llms.MessageContent {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, "Hello Rick, introduce yourself."),
		llms.TextParts(llms.ChatMessageTypeAI, "Hi! I'm Rick, your AI assistant. How can I help you today?"),
	}

	if initialPrompt != "" {
		content = append(content,
			llms.TextParts(llms.ChatMessageTypeHuman, initialPrompt),
			llms.TextParts(llms.ChatMessageTypeAI, "I understand. How can I help you?"),
		)
	}
	return content
}

// HandlePerplexityResponse manages Perplexity-specific response handling
func HandlePerplexityResponse(content *[]llms.MessageContent, response string) {
	// Keep conversation history manageable for Perplexity
	if len(*content) > 6 {
		*content = (*content)[len(*content)-6:]
	}
	*content = append(*content, llms.TextParts(llms.ChatMessageTypeAI, response))
}
