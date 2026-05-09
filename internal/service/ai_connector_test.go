package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAIConnector_MockAnalysis(t *testing.T) {
	connector := &AIConnector{
		apiURL: "",
		apiKey: "",
		model:  "test-model",
	}

	result, err := connector.Analyze("test document text")
	assert.NoError(t, err)

	var analysis AIAnalysisResult
	err = json.Unmarshal(result, &analysis)
	assert.NoError(t, err)
	assert.NotEmpty(t, analysis.Risks)
	assert.NotEmpty(t, analysis.Recommendations)
	assert.Contains(t, analysis.Summary, "тестовый режим")
}

func TestAIConnector_ModelName(t *testing.T) {
	connector := &AIConnector{model: "deepseek-chat"}
	assert.Equal(t, "deepseek-chat", connector.ModelName())

	connector.model = "gpt-4"
	assert.Equal(t, "gpt-4", connector.ModelName())
}

func TestAIConnector_Analyze_ValidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		response := `{
			"choices": [{
				"message": {
					"content": "{\"risks\":[\"risk1\"],\"fraud_indicators\":[],\"recommendations\":[\"rec1\"],\"summary\":\"test\"}"
				}
			}],
			"usage": {"total_tokens": 100}
		}`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	connector := &AIConnector{
		apiURL: server.URL,
		apiKey: "test-key",
		model:  "test-model",
		client: server.Client(),
	}

	result, err := connector.Analyze("test text")
	assert.NoError(t, err)

	var analysis AIAnalysisResult
	err = json.Unmarshal(result, &analysis)
	assert.NoError(t, err)
	assert.Contains(t, analysis.Risks, "risk1")
}

func TestAIConnector_Analyze_InvalidJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"choices": [{"message": {"content": "not a json"}}]}`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	connector := &AIConnector{
		apiURL: server.URL,
		apiKey: "test-key",
		model:  "test-model",
		client: server.Client(),
	}

	result, err := connector.Analyze("test text")
	assert.NoError(t, err)

	var analysis AIAnalysisResult
	err = json.Unmarshal(result, &analysis)
	assert.NoError(t, err)
	assert.Equal(t, "not a json", analysis.Summary)
}

func TestAIConnector_Analyze_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	connector := &AIConnector{
		apiURL: server.URL,
		apiKey: "test-key",
		model:  "test-model",
		client: server.Client(),
	}

	_, err := connector.Analyze("test text")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ai api error 500")
}

func TestAIConnector_Analyze_EmptyAPIURL(t *testing.T) {
	connector := &AIConnector{
		apiURL: "",
		apiKey: "test-key",
		model:  "test-model",
	}

	result, err := connector.Analyze("test text")
	assert.NoError(t, err)
	assert.Contains(t, string(result), "тестовый режим")
}

func TestAIConnector_Analyze_EmptyAPIKey(t *testing.T) {
	connector := &AIConnector{
		apiURL: "https://api.test.com",
		apiKey: "",
		model:  "test-model",
	}

	result, err := connector.Analyze("test text")
	assert.NoError(t, err)
	assert.Contains(t, string(result), "тестовый режим")
}

func TestAIConnector_Analyze_NetworkTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	defer server.Close()

	connector := &AIConnector{
		apiURL: server.URL,
		apiKey: "test-key",
		model:  "test-model",
		client: &http.Client{Timeout: 100 * time.Millisecond},
	}

	_, err := connector.Analyze("test text")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestAIConnector_Analyze_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"choices": []}`))
	}))
	defer server.Close()

	connector := &AIConnector{
		apiURL: server.URL,
		apiKey: "test-key",
		model:  "test-model",
		client: server.Client(),
	}

	_, err := connector.Analyze("test text")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

func TestAIConnector_ValidateJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		valid    bool
	}{
		{"valid json", `{"risks":[]}`, true},
		{"invalid json", `not json`, false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := json.Valid([]byte(tt.response))
			assert.Equal(t, tt.valid, valid)
		})
	}
}

func TestNewAIConnector(t *testing.T) {
	origURL := os.Getenv("AI_API_URL")
	origKey := os.Getenv("AI_API_KEY")
	origModel := os.Getenv("AI_MODEL")
	defer func() {
		os.Setenv("AI_API_URL", origURL)
		os.Setenv("AI_API_KEY", origKey)
		os.Setenv("AI_MODEL", origModel)
	}()

	os.Setenv("AI_API_URL", "")
	os.Setenv("AI_API_KEY", "")
	os.Setenv("AI_MODEL", "")
	connector := NewAIConnector()
	assert.NotNil(t, connector)
	assert.Equal(t, "deepseek-chat", connector.ModelName())

	os.Setenv("AI_API_URL", "https://custom.api.com")
	os.Setenv("AI_API_KEY", "custom-key")
	os.Setenv("AI_MODEL", "gpt-4")
	connector2 := NewAIConnector()
	assert.NotNil(t, connector2)
	assert.Equal(t, "gpt-4", connector2.ModelName())
}