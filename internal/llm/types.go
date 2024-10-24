package llm

type Provider interface {
	GenerateResponse(prompt string) (string, error)
	// Other common LLM methods
}
