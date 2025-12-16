package entities

// Request para refinar texto com IA
type RefineTextRequest struct {
	Text   string `json:"text"`
	Action string `json:"action"` // "summarize", "refine", "expand"
}

// Response da IA
type RefineTextResponse struct {
	OriginalText string `json:"original_text"`
	RefinedText  string `json:"refined_text"`
	Action       string `json:"action"`
	Model        string `json:"model"`
}

// Tipos de ações
const (
	ActionSummarize = "summarize" // Resumir
	ActionRefine    = "refine"    // Melhorar/Refinar
	ActionExpand    = "expand"    // Expandir/Detalhar
)
