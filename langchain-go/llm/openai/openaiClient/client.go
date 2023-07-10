package openaiClient

import (
	"context"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"net/http"
	"time"
)

type OpenAiClient struct {
	*llmSchema.BaseAIClient
	APIKey          string
	OrganizationKey string
}

func (c *OpenAiClient) addHeaders(req *http.Request) {
	// Call the parent method to add default headers
	c.BaseAIClient.AddHeaders(req)

	// custom headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	if c.OrganizationKey != "" {
		req.Header.Set("OpenAI-Organization", c.OrganizationKey)
	}

}

func (c *OpenAiClient) Create(prompts []string, input map[string]interface{}) ([]CompletionResponsePayload, error) {
	var err error
	var response []CompletionResponsePayload

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	requestPayload, err := c.createCompletionRequestPayload(input)
	for _, prompt := range prompts {
		requestPayload.Prompt = prompt
		if err != nil {
			return nil, err
		}
		newResponse, err := c.BaseAIClient.Create(ctx, requestPayload)
		response = append(response, newResponse.(CompletionResponsePayload))
		if err != nil {
			return nil, err
		}
	}
	return response, err
}

func (c *OpenAiClient) createCompletionResponsePayload() llmSchema.ResponsePayload {
	return NewCompletionResponsePayload()
}

func (c *OpenAiClient) createCompletionRequestPayload(input map[string]interface{}) (*CompletionRequestPayload, error) {
	payload := &CompletionRequestPayload{}

	for key, value := range input {
		switch key {
		case "model_name":
			if s, ok := value.(string); ok && s != "" {
				payload.Model = s
			}
		case "prompt":
			if p, ok := value.(string); ok && len(p) > 0 {
				payload.Prompt = p
			}
		case "temperature":
			if f, ok := value.(float64); ok && f != 0 {
				payload.Temperature = f
			}
		case "max_tokens":
			if i, ok := value.(int); ok && i != 0 {
				payload.MaxTokens = i
			}
		case "top_p":
			if f, ok := value.(float64); ok && f != 0 {
				payload.TopP = f
			}
		case "frequency_penalty":
			if f, ok := value.(float64); ok && f != 0 {
				payload.FrequencyPenalty = f
			}
		case "presence_penalty":
			if f, ok := value.(float64); ok && f != 0 {
				payload.PresencePenalty = f
			}
		case "n":
			if i, ok := value.(int); ok && i != 0 {
				payload.N = i
			}
		case "best_of":
			if i, ok := value.(int); ok && i != 0 {
				payload.BestOf = i
			}
		case "model_kwargs":
			if m, ok := value.(map[string]interface{}); ok && len(m) > 0 {
				kwargs := make(map[string]string)
				for k, v := range m {
					if vs, ok := v.(string); ok && vs != "" {
						kwargs[k] = vs
					}
				}
				payload.ModelKwargs = kwargs
			}
		case "batch_size":
			if i, ok := value.(int); ok && i != 0 {
				payload.BatchSize = i
			}
		case "logit_bias":
			if m, ok := value.(map[string]interface{}); ok && len(m) > 0 {
				bias := make(map[string]float64)
				for k, v := range m {
					if vf, ok := v.(float64); ok && vf != 0 {
						bias[k] = vf
					}
				}
				payload.LogitBias = bias
			}
		case "max_retries":
			if i, ok := value.(int); ok && i != 0 {
				payload.MaxRetries = i
			}
		case "streaming":
			if b, ok := value.(bool); ok {
				payload.Streaming = b
			}
		case "stop":
			if a, ok := value.([]interface{}); ok && len(a) > 0 {
				stopWords := make([]string, 0, len(a))
				for _, v := range a {
					if s, ok := v.(string); ok && s != "" {
						stopWords = append(stopWords, s)
					}
				}
				payload.StopWords = stopWords
			}
		default:
			// ignore unknown fields
		}
	}

	return payload, nil
}

func NewOpenAiClient(APIKey string, APIOrganization string, APIBaseURL string, maxRetries int) (OpenAiClient, error) {
	var err error
	// initialize the base openaiClient
	// ( apiBaseURL string, maxRetries int, responsePayload ResponsePayload
	baseAiClient := llmSchema.NewBaseAIClient(APIBaseURL, maxRetries, NewCompletionResponsePayload())

	// initialize the OpenAI openaiClient
	client := OpenAiClient{
		baseAiClient,
		APIKey,
		APIOrganization,
	}

	return client, err
}
