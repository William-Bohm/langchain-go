package callbacks

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/tools"
)

type StdOutCallbackHandler struct {
	callbackSchema.BaseCallbackHandler
	Color *string
}

func NewStdOutCallbackHandler(color *string) *StdOutCallbackHandler {
	return &StdOutCallbackHandler{
		Color: color,
	}
}

func (h *StdOutCallbackHandler) OnLlmStart(serialized map[string]interface{}, prompts []string, kwargs map[string]interface{}) {
}

func (h *StdOutCallbackHandler) OnLlmEnd(response llmSchema.LLMResult, kwargs map[string]interface{}) {
}

func (h *StdOutCallbackHandler) OnLlmNewToken(token string, kwargs map[string]interface{}) {}

func (h *StdOutCallbackHandler) OnLlmError(err error, kwargs map[string]interface{}) {}

func (h *StdOutCallbackHandler) OnChainStart(serialized map[string]interface{}, inputs map[string]interface{}, kwargs map[string]interface{}) {
	className := serialized["name"].(string)
	fmt.Printf("\n\n\033[1m> Entering new %s chain...\033[0m", className)
}

func (h *StdOutCallbackHandler) OnChainEnd(outputs map[string]interface{}, kwargs map[string]interface{}) {
	fmt.Println("\n\033[1m> Finished chain.\033[0m")
}

func (h *StdOutCallbackHandler) OnChainError(err error, kwargs map[string]interface{}) {}

func (h *StdOutCallbackHandler) OnToolStart(serialized map[string]interface{}, inputStr string, kwargs map[string]interface{}) {
}

func (h *StdOutCallbackHandler) OnAgentAction(action callbackSchema.AgentAction, color *string, kwargs map[string]interface{}) {
	printColor := h.Color
	if color != nil {
		printColor = color
	}
	tools.PrintText(action.Log, printColor, "")
}

func (h *StdOutCallbackHandler) OnToolEnd(output string, color *string, observationPrefix *string, llmPrefix *string, kwargs map[string]interface{}) {
	if observationPrefix != nil {
		tools.PrintText("\n"+*observationPrefix, h.Color, "")
	}
	printColor := h.Color
	if color != nil {
		printColor = color
	}
	tools.PrintText(output, printColor, "")
	if llmPrefix != nil {
		tools.PrintText("\n"+*llmPrefix, h.Color, "")
	}
}

func (h *StdOutCallbackHandler) OnToolError(err error, kwargs map[string]interface{}) {}

func (h *StdOutCallbackHandler) OnText(text string, color *string, end string, kwargs map[string]interface{}) {
	printColor := h.Color
	if color != nil {
		printColor = color
	}
	tools.PrintText(text, printColor, end)
}

func (h *StdOutCallbackHandler) OnAgentFinish(finish callbackSchema.AgentFinish, color *string, kwargs map[string]interface{}) {
	printColor := h.Color
	if color != nil {
		printColor = color
	}
	tools.PrintText(finish.Log, printColor, "\n")
}
