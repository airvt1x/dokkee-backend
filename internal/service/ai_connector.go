package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type AIAnalysisResult struct {
	Risks          []string `json:"risks"`
	FraudIndicators []string `json:"fraud_indicators"`
	Recommendations []string `json:"recommendations"`
	Summary        string   `json:"summary"`
}

type AIConnector struct {
	apiURL string
	apiKey string
	client *http.Client
}

func NewAIConnector() *AIConnector {
	return &AIConnector{
		apiURL: os.Getenv("AI_API_URL"),
		apiKey: os.Getenv("AI_API_KEY"),
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *AIConnector) Analyze(anonymizedText string) ([]byte, error) {
	if c.apiURL == "" {
		return c.mockAnalysis(anonymizedText)
	}

	prompt := fmt.Sprintf(`Проанализируй следующий юридический документ и выяви:
1. Юридические и финансовые риски
2. Признаки мошеннических схем
3. Несоответствия нормативным требованиям
4. Рекомендации

Документ (персональные данные обезличены):
%s

Верни ответ в формате JSON со следующими полями:
- risks: список рисков
- fraud_indicators: признаки мошенничества
- recommendations: рекомендации
- summary: краткое резюме`, anonymizedText)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"response_format": map[string]string{"type": "json_object"},
	})

	req, err := http.NewRequest(http.MethodPost, c.apiURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ai request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ai api error %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("empty response from ai")
	}

	return []byte(apiResp.Choices[0].Message.Content), nil
}

func (c *AIConnector) mockAnalysis(text string) ([]byte, error) {
	result := AIAnalysisResult{
		Risks:           []string{"Документ требует проверки юристом"},
		FraudIndicators: []string{},
		Recommendations: []string{"Рекомендуется детальная проверка условий договора"},
		Summary:         fmt.Sprintf("Документ проанализирован (тестовый режим). Длина текста: %d символов.", len(text)),
	}
	return json.Marshal(result)
}
