package callbacks

// This file is in the original Python version, I dont understand its purpose
// perhaps it was not fully implemented yet...

import (
	"fmt"
	"os"
)

type BaseCallbackHandler interface {
}

type StreamingStdOutCallbackHandler struct {
}

func (s *StreamingStdOutCallbackHandler) OnLLMStart(serialized map[string]interface{}, prompts []string, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnLLMNewToken(token string, kwargs map[string]interface{}) {
	fmt.Fprint(os.Stdout, token)
}

func (s *StreamingStdOutCallbackHandler) OnLLMEnd(response interface{}, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnLLMError(err interface{}, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnChainStart(serialized map[string]interface{}, inputs map[string]interface{}, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnChainEnd(outputs map[string]interface{}, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnChainError(err interface{}, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnToolStart(serialized map[string]interface{}, inputStr string, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnAgentAction(action interface{}, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnToolEnd(output string, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnToolError(err interface{}, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnText(text string, kwargs map[string]interface{}) {
}

func (s *StreamingStdOutCallbackHandler) OnAgentFinish(finish interface{}, kwargs map[string]interface{}) {
}
