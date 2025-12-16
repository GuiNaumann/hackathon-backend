package repository_impl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/settings_loader"
	"net/http"
	"time"
)

type AIRepositoryImpl struct {
	settings *settings_loader.SettingsLoader
	client   *http.Client
}

func NewAIRepositoryImpl(settings *settings_loader.SettingsLoader) *AIRepositoryImpl {
	timeout := time.Duration(settings.GetAIRequestTimeout()) * time.Second

	return &AIRepositoryImpl{
		settings: settings,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Estruturas para comunicação com a API Gemini
type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	Temperature float64 `json:"temperature"`
	MaxOutputTokens int `json:"maxOutputTokens"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

func (r *AIRepositoryImpl) ProcessText(ctx context.Context, req *entities.AIRequest) (*entities.AIResponse, error) {
	// Verificar se IA está habilitada
	if !r.settings.IsAIEnabled() {
		return nil, errors.New("funcionalidade de IA não está habilitada. Configure ai.gemini_api_key no settings.toml")
	}

	// Validações
	if req.Text == "" {
		return nil, errors.New("texto não pode estar vazio")
	}

	if req.Prompt == "" {
		return nil, errors.New("prompt não pode estar vazio")
	}

	if len(req.Text) > 10000 {
		return nil, errors.New("texto muito longo (máximo 10000 caracteres)")
	}

	// Combinar prompt com texto
	fullPrompt := fmt.Sprintf("%s\n\nTexto:\n%s", req.Prompt, req.Text)

	// Montar requisição para Gemini
	geminiReq := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: fullPrompt},
				},
			},
		},
		GenerationConfig: geminiGenerationConfig{
			Temperature: r.settings.GetAITemperature(),
			MaxOutputTokens: r.settings.GetAIMaxTokens(),
		},
	}

	payload, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao montar payload: %w", err)
	}

	// Criar requisição HTTP
	apiKey := r.settings.GetGeminiAPIKey()
	model := r.settings.GetGeminiModel()
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", model)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	httpReq.Header.Set("X-Goog-Api-Key", apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Executar requisição
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar API Gemini: %w", err)
	}
	defer resp.Body.Close()

	// Parse da resposta
	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Verificar erro na resposta
	if geminiResp.Error != nil {
		return nil, fmt.Errorf("erro da API Gemini: %s", geminiResp.Error.Message)
	}

	// Verificar se tem candidates
	if len(geminiResp.Candidates) == 0 {
		return nil, errors.New("resposta vazia da API Gemini")
	}

	// Verificar se tem parts
	if len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("resposta sem conteúdo da API Gemini")
	}

	generatedText := geminiResp.Candidates[0].Content.Parts[0].Text

	return &entities.AIResponse{
		OriginalText: req.Text,
		GeneratedText: generatedText,
		Prompt:       req.Prompt,
		Model:        model,
	}, nil
}
