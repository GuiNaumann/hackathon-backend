package entities

// Request para processar texto com IA
type AIRequest struct {
	Text   string `json:"text"`
	Prompt string `json:"prompt"`
}

// Response da IA
type AIResponse struct {
	OriginalText string `json:"original_text"`
	GeneratedText string `json:"generated_text"`
	Prompt       string `json:"prompt"`
	Model        string `json:"model"`
}
