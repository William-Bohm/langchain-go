package openaiClient

import (
	"encoding/json"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
)

// optional openai endpoint settings
type CompletionRequestPayload struct {
	Model              string             `json:"model_name"`
	Prompt             string             `json:"prompt"`
	Temperature        float64            `json:"temperature,omitempty"`
	MaxTokens          int                `json:"max_tokens,omitempty"`
	TopP               float64            `json:"top_p,omitempty"`
	FrequencyPenalty   float64            `json:"frequency_penalty,omitempty"`
	PresencePenalty    float64            `json:"presence_penalty,omitempty"`
	N                  int                `json:"n,omitempty"`
	BestOf             int                `json:"best_of,omitempty"`
	ModelKwargs        map[string]string  `json:"model_kwargs,omitempty"`
	OpenaiApiKey       string             `json:"-"`
	OpenaiApiBase      string             `json:"-"`
	OpenaiOrganization string             `json:"-"`
	BatchSize          int                `json:"batch_size,omitempty"`
	RequestTimeout     float64            `json:"-"`
	LogitBias          map[string]float64 `json:"logit_bias,omitempty"`
	MaxRetries         int                `json:"max_retries,omitempty"`
	Streaming          bool               `json:"streaming,omitempty"`
	StopWords          []string           `json:"stop,omitempty"`
}

func (p *CompletionRequestPayload) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// open ai endpoint response
type CompletionResponsePayload struct {
	llmSchema.ResponsePayload
	ID      string  `json:"id,omitempty"`
	Object  string  `json:"object,omitempty"`
	Created float64 `json:"created,omitempty"`
	Model   string  `json:"model,omitempty"`

	Usage struct {
		CompletionTokens float64 `json:"completion_tokens,omitempty"`
		PromptTokens     float64 `json:"prompt_tokens,omitempty"`
		TotalTokens      float64 `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`

	Choices []struct {
		FinishReason string      `json:"finish_reason,omitempty"`
		Index        float64     `json:"index,omitempty"`
		Logprobs     interface{} `json:"logprobs,omitempty"`
		Text         string      `json:"text,omitempty"`
	} `json:"choices,omitempty"`
}

func (p CompletionResponsePayload) FromJSON(data []byte) (llmSchema.ResponsePayload, error) {
	var response CompletionResponsePayload
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
func (p CompletionResponsePayload) NewResponsePayload() llmSchema.ResponsePayload {
	return CompletionResponsePayload{}
}

func NewCompletionResponsePayload() llmSchema.ResponsePayload {
	return CompletionResponsePayload{}
}

type Messages struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
