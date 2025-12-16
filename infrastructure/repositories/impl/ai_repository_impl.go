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

// Estruturas para comunicação com a API Groq
type groqRequest struct {
	Model       string        `json:"model"`
	Messages    []groqMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqResponse struct {
	Choices []struct {
		Message groqMessage `json:"message"`
	} `json:"choices"`
	Model string `json:"model"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (r *AIRepositoryImpl) RefineText(ctx context.Context, req *entities.RefineTextRequest) (*entities.RefineTextResponse, error) {
	// Verificar se IA está habilitada
	if !r.settings.IsAIEnabled() {
		return nil, errors.New("funcionalidade de IA não está habilitada.  Configure ai.groq_api_key no settings.toml")
	}

	// Validações
	if req.Text == "" {
		return nil, errors.New("texto não pode estar vazio")
	}

	if len(req.Text) > 5000 {
		return nil, errors.New("texto muito longo (máximo 5000 caracteres)")
	}

	// Definir prompt baseado na ação
	systemPrompt, userPrompt := r.buildPrompts(req.Action, req.Text)

	// Montar requisição para Groq
	groqReq := groqRequest{
		Model: r.settings.GetGroqModel(),
		Messages: []groqMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: r.settings.GetAITemperature(),
		MaxTokens:   r.settings.GetAIMaxTokens(),
	}

	payload, err := json.Marshal(groqReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao montar payload: %w", err)
	}

	// Criar requisição HTTP
	apiKey := r.settings.GetGroqAPIKey()

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Executar requisição
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar API Groq: %w", err)
	}
	defer resp.Body.Close()

	// Parse da resposta
	var groqResp groqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Verificar erro na resposta
	if groqResp.Error != nil {
		return nil, fmt.Errorf("erro da API Groq: %s", groqResp.Error.Message)
	}

	// Verificar se tem choices
	if len(groqResp.Choices) == 0 {
		return nil, errors.New("resposta vazia da API Groq")
	}

	refinedText := groqResp.Choices[0].Message.Content

	return &entities.RefineTextResponse{
		OriginalText: req.Text,
		RefinedText:  refinedText,
		Action:       req.Action,
		Model:        groqResp.Model,
	}, nil
}

func (r *AIRepositoryImpl) buildPrompts(action, text string) (systemPrompt, userPrompt string) {
	switch action {
	case entities.ActionSummarize:
		systemPrompt = "Você é um assistente especializado em criar resumos claros e concisos de textos, mantendo as informações mais importantes."
		userPrompt = fmt.Sprintf("Resuma o seguinte texto de forma clara e objetiva:\n\n%s", text)

	case entities.ActionExpand:
		systemPrompt = "Você é um assistente que expande e detalha textos, adicionando mais contexto e informações relevantes de forma profissional."
		userPrompt = fmt.Sprintf("Expanda e detalhe o seguinte texto, mantendo o tom profissional:\n\n%s", text)

	case entities.ActionRefine:
		fallthrough
	default:
		systemPrompt = "Você é um assistente que melhora textos tornando-os mais claros, profissionais e bem estruturados, mantendo a ideia original."
		userPrompt = fmt.Sprintf("Melhore e refine o seguinte texto, tornando-o mais claro e profissional:\n\n%s", text)
	}

	return systemPrompt, userPrompt
}
