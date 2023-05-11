package callbacks

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"strings"
)

func getOpenAIModelCostPer1kTokens(modelName string, isCompletion bool) (float64, error) {
	modelCostMapping := map[string]float64{
		"gpt-4":                     0.03,
		"gpt-4-0314":                0.03,
		"gpt-4-completion":          0.06,
		"gpt-4-0314-completion":     0.06,
		"gpt-4-32k":                 0.06,
		"gpt-4-32k-0314":            0.06,
		"gpt-4-32k-completion":      0.12,
		"gpt-4-32k-0314-completion": 0.12,
		"gpt-3.5-turbo":             0.002,
		"gpt-3.5-turbo-0301":        0.002,
		"text-ada-001":              0.0004,
		"ada":                       0.0004,
		"text-babbage-001":          0.0005,
		"babbage":                   0.0005,
		"text-curie-001":            0.002,
		"curie":                     0.002,
		"text-davinci-003":          0.02,
		"text-davinci-002":          0.02,
		"code-davinci-002":          0.02,
	}

	lowerModelName := strings.ToLower(modelName)
	if isCompletion && strings.HasPrefix(lowerModelName, "gpt-4") {
		lowerModelName += "-completion"
	}
	cost, ok := modelCostMapping[lowerModelName]
	if !ok {
		keys := make([]string, 0, len(modelCostMapping))
		for k := range modelCostMapping {
			keys = append(keys, k)
		}
		return 0, fmt.Errorf("unknown model: %s. Please provide a valid OpenAI model name. Known models are: %s",
			modelName, strings.Join(keys, ", "))
	}
	return cost, nil
}

type OpenAICallbackHandler struct {
	TotalTokens        int
	PromptTokens       int
	CompletionTokens   int
	SuccessfulRequests int
	TotalCost          float64
}

func (o *OpenAICallbackHandler) String() string {
	return fmt.Sprintf("Tokens Used: %d\n\tPrompt Tokens: %d\n\tCompletion Tokens: %d\nSuccessful Requests: %d\nTotal Cost (USD): $%.2f",
		o.TotalTokens, o.PromptTokens, o.CompletionTokens, o.SuccessfulRequests, o.TotalCost)
}

func (o *OpenAICallbackHandler) AlwaysVerbose() bool {
	return true
}

func (o *OpenAICallbackHandler) OnLLMStart(serialized map[string]interface{}, prompts []string, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnLLMNewToken(token string, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnLLMEnd(response llmSchema.LLMResult, kwargs map[string]interface{}) {
	if response.LLMOutput != nil {
		o.SuccessfulRequests += 1
		if tokenUsage, ok := response.LLMOutput["token_usage"].(map[string]interface{}); ok {
			if modelName, ok := response.LLMOutput["model_name"].(string); ok {
				completionCost, err := getOpenAIModelCostPer1kTokens(modelName, true)
				if err != nil {
					println("failed to get model tokens")
				}
				promptCost, err := getOpenAIModelCostPer1kTokens(modelName, false)

				if completionTokens, ok := tokenUsage["completion_tokens"].(float64); ok {
					completionCost *= completionTokens / 1000
				}
				if err != nil {
					println("failed to get model tokens")
				}

				if promptTokens, ok := tokenUsage["prompt_tokens"].(float64); ok {
					promptCost *= promptTokens / 1000
				}

				o.TotalCost += promptCost + completionCost

			}
			if totalTokens, ok := tokenUsage["total_tokens"].(float64); ok {
				o.TotalTokens += int(totalTokens)
			}
			if promptTokens, ok := tokenUsage["prompt_tokens"].(float64); ok {
				o.PromptTokens += int(promptTokens)
			}
			if completionTokens, ok := tokenUsage["completion_tokens"].(float64); ok {
				o.CompletionTokens += int(completionTokens)
			}
		}
	}
}

func (o *OpenAICallbackHandler) OnLLMError(err error, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnChainStart(serialized map[string]interface{}, inputs map[string]interface{}, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnChainEnd(outputs map[string]interface{}, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnChainError(err error, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnToolStart(serialized map[string]interface{}, inputStr string, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnToolEnd(output string, color string, observationPrefix string, llmPrefix string, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnToolError(err error, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnText(text string, color string, end string, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnAgentAction(action callbackSchema.AgentAction, kwargs map[string]interface{}) {
}

func (o *OpenAICallbackHandler) OnAgentFinish(finish callbackSchema.AgentFinish, color string, kwargs map[string]interface{}) {
}
