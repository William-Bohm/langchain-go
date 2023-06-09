package callbacks

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/agent/agentSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/tools"
)

type StdOutCallbackHandler struct {
	callbackSchema.BaseCallbackHandler
	Color string
}

func NewStdOutCallbackHandler(color string) *StdOutCallbackHandler {
	return &StdOutCallbackHandler{
		Color: color,
	}
}

func (h *StdOutCallbackHandler) OnLlmStart(serialized map[string]interface{}, prompts []string, verbose bool, args ...interface{}) {
}

func (h *StdOutCallbackHandler) OnLlmEnd(response llmSchema.LLMResult, verbose bool, args ...interface{}) {
}

func (h *StdOutCallbackHandler) OnLlmNewToken(token string, verbose bool, args ...interface{}) {}

func (h *StdOutCallbackHandler) OnLlmError(err error, verbose bool, args ...interface{}) {}

func (h *StdOutCallbackHandler) OnChainStart(serialized map[string]interface{}, inputs map[string]interface{}, verbose bool, args ...interface{}) {
	className := serialized["name"].(string)
	fmt.Printf("\n\n\033[1m> Entering new %s chain...\033[0m", className)
}

func (h *StdOutCallbackHandler) OnChainEnd(outputs map[string]interface{}, verbose bool, args ...interface{}) {
	fmt.Println("\n\033[1m> Finished chain.\033[0m")
}

func (h *StdOutCallbackHandler) OnChainError(err error, verbose bool, args ...interface{}) {}

func (h *StdOutCallbackHandler) OnToolStart(serialized map[string]interface{}, inputStr string, verbose bool, args ...interface{}) {
}

func (h *StdOutCallbackHandler) OnAgentAction(action agentSchema.AgentAction, color string, verbose bool, args ...interface{}) {
	printColor := h.Color
	if color != "" {
		printColor = color
	}
	tools.PrintText(action.Log, printColor, "")
}

func (h *StdOutCallbackHandler) OnToolEnd(output string, color string, observationPrefix *string, llmPrefix *string, verbose bool, args ...interface{}) {
	if observationPrefix != nil {
		tools.PrintText("\n"+*observationPrefix, h.Color, "")
	}
	printColor := h.Color
	if color != "" {
		printColor = color
	}
	tools.PrintText(output, printColor, "")
	if llmPrefix != nil {
		tools.PrintText("\n"+*llmPrefix, h.Color, "")
	}
}

func (h *StdOutCallbackHandler) OnToolError(err error, verbose bool, args ...interface{}) {}

func (h *StdOutCallbackHandler) OnText(text string, color string, end string, verbose bool, args ...interface{}) {
	printColor := h.Color
	if color != "" {
		printColor = color
	}
	tools.PrintText(text, printColor, end)
}

func (h *StdOutCallbackHandler) OnAgentFinish(finish agentSchema.AgentFinish, color string, verbose bool, args ...interface{}) {
	printColor := h.Color
	if color != "" {
		printColor = color
	}
	tools.PrintText(finish.Log, printColor, "\n")
}
