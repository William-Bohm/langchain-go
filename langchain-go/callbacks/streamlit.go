package callbacks

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/agent/agentSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
)

type StreamlitCallbackHandler struct {
}

func (s *StreamlitCallbackHandler) OnLLMStart(serialized map[string]interface{}, prompts []string, kwargs map[string]interface{}) {
	fmt.Println("Prompts after formatting:")
	for _, prompt := range prompts {
		fmt.Println(prompt)
	}
}

func (s *StreamlitCallbackHandler) OnLLMNewToken(token string, kwargs map[string]interface{}) {
	// Do nothing.
}

func (s *StreamlitCallbackHandler) OnLLMEnd(response llmSchema.LLMResult, kwargs map[string]interface{}) {
	// Do nothing.
}

func (s *StreamlitCallbackHandler) OnLLMError(err interface{}, kwargs map[string]interface{}) {
	// Do nothing.
}

func (s *StreamlitCallbackHandler) OnChainStart(serialized map[string]interface{}, inputs map[string]interface{}, kwargs map[string]interface{}) {
	className := serialized["name"]
	fmt.Printf("Entering new %s chain...\n", className)
}

func (s *StreamlitCallbackHandler) OnChainEnd(outputs map[string]interface{}, kwargs map[string]interface{}) {
	fmt.Println("Finished chain.")
}

func (s *StreamlitCallbackHandler) OnChainError(err interface{}, kwargs map[string]interface{}) {
	// Do nothing.
}

func (s *StreamlitCallbackHandler) OnToolStart(serialized map[string]interface{}, inputStr string, kwargs map[string]interface{}) {
	// Do nothing.
}

func (s *StreamlitCallbackHandler) OnAgentAction(action agentSchema.AgentAction, kwargs map[string]interface{}) {
	fmt.Println(action.Log)
}

func (s *StreamlitCallbackHandler) OnToolEnd(output string, observationPrefix string, llmPrefix string, kwargs map[string]interface{}) {
	fmt.Printf("%s%s\n", observationPrefix, output)
	fmt.Println(llmPrefix)
}

func (s *StreamlitCallbackHandler) OnToolError(err interface{}, kwargs map[string]interface{}) {
	// Do nothing.
}

func (s *StreamlitCallbackHandler) OnText(text string, kwargs map[string]interface{}) {
	fmt.Println(text)
}

func (s *StreamlitCallbackHandler) OnAgentFinish(finish agentSchema.AgentFinish, kwargs map[string]interface{}) {
	fmt.Println(finish.Log)
}
