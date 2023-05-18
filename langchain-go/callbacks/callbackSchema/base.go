package callbackSchema

import (
	"github.com/William-Bohm/langchain-go/langchain-go/agent/agentSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
)

type BaseCallbackHandler interface {
	AlwaysVerbose() bool
	IgnoreLLM() bool
	IgnoreChain() bool
	IgnoreAgent() bool
	OnLLMStart(serialized map[string]interface{}, prompts []string, args ...interface{}) (interface{}, error)
	OnLLMNewToken(token string, args ...interface{}) (interface{}, error)
	OnLLMEnd(response llmSchema.LLMResult, args ...interface{}) (interface{}, error)
	OnLLMError(err error, args ...interface{}) (interface{}, error)
	OnChainStart(serialized map[string]interface{}, inputs map[string]interface{}, args ...interface{}) (interface{}, error)
	OnChainEnd(outputs map[string]interface{}, args ...interface{}) (interface{}, error)
	OnChainError(err error, args ...interface{}) (interface{}, error)
	OnToolStart(serialized map[string]interface{}, inputStr string, args ...interface{}) (interface{}, error)
	OnToolEnd(output string, args ...interface{}) (interface{}, error)
	OnToolError(err error, args ...interface{}) (interface{}, error)
	OnText(text string, args ...interface{}) (interface{}, error)
	OnAgentAction(action agentSchema.AgentAction, args ...interface{}) (interface{}, error)
	OnAgentFinish(finish agentSchema.AgentFinish, args ...interface{}) (interface{}, error)
}

type BaseCallbackManager interface {
	BaseCallbackHandler
	AddHandler(callback BaseCallbackHandler)
	RemoveHandler(handler BaseCallbackHandler)
	SetHandler(handler BaseCallbackHandler)
	SetHandlers(handlers []BaseCallbackHandler)
}

type CallbackManager struct {
	handlers []BaseCallbackHandler
}

func NewCallbackManager(handlers []BaseCallbackHandler) *CallbackManager {
	return &CallbackManager{handlers: handlers}
}

func (c *CallbackManager) OnLLMStart(serialized map[string]interface{}, prompts []string, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreLLM() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnLLMStart(serialized, prompts, args...)
			}
		}
	}
}

func (c *CallbackManager) OnLLMNewToken(token string, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreLLM() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnLLMNewToken(token, args)
			}
		}
	}
}

func (c *CallbackManager) OnLLMEnd(response llmSchema.LLMResult, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreLLM() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnLLMEnd(response, args)
			}
		}
	}
}

func (c *CallbackManager) OnLLMError(err error, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreLLM() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnLLMError(err, args)
			}
		}
	}
}

func (c *CallbackManager) OnChainStart(serialized map[string]interface{}, inputs map[string]interface{}, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreChain() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnChainStart(serialized, inputs, kwargs)
			}
		}
	}
}

func (c *CallbackManager) OnChainEnd(outputs map[string]interface{}, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreChain() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnChainEnd(outputs, args)
			}
		}
	}
}

func (c *CallbackManager) OnChainError(err error, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreChain() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnChainError(err, args)
			}
		}
	}
}

func (c *CallbackManager) OnToolStart(serialized map[string]interface{}, input_str string, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreAgent() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnToolStart(serialized, input_str, args)
			}
		}
	}
}

func (c *CallbackManager) OnAgentAction(action agentSchema.AgentAction, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreAgent() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnAgentAction(action, args)
			}
		}
	}
}

func (c *CallbackManager) OnToolEnd(output string, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreAgent() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnToolEnd(output, args)
			}
		}
	}
}

func (c *CallbackManager) OnToolError(err error, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreAgent() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnToolError(err, args)
			}
		}
	}
}

func (c *CallbackManager) OnText(text string, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if verbose || handler.AlwaysVerbose() {
			handler.OnText(text, args)
		}
	}
}

func (c *CallbackManager) OnAgentFinish(finish agentSchema.AgentFinish, verbose bool, args ...interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreAgent() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnAgentFinish(finish, args)
			}
		}
	}
}

func (c *CallbackManager) AddHandler(handler BaseCallbackHandler) {
	c.handlers = append(c.handlers, handler)
}

func (c *CallbackManager) RemoveHandler(handler BaseCallbackHandler) {
	for i, h := range c.handlers {
		if h == handler {
			c.handlers = append(c.handlers[:i], c.handlers[i+1:]...)
			break
		}
	}
}

func (c *CallbackManager) SetHandlers(handlers []BaseCallbackHandler) {
	c.handlers = handlers
}

func handleEventForHandler(
	handler BaseCallbackHandler,
	eventName string,
	ignoreConditionName *string,
	verbose bool,
	args ...interface{}) {

	if ignoreConditionName == nil || !handlerIgnoreCondition(handler, *ignoreConditionName) {
		if verbose || handler.AlwaysVerbose() {
			go handleEvent(handler, eventName, args...)
		}
	}
}

// TODO: add all case's
func handlerIgnoreCondition(handler BaseCallbackHandler, conditionName string) bool {
	switch conditionName {
	case "ignore_llm":
		return handler.IgnoreLLM()
	case "ignore_chain":
		return handler.IgnoreChain()
	case "ignore_agent":
		return handler.IgnoreAgent()
	default:
		return false
	}
}

// TODO: add all case's
func handleEvent(handler BaseCallbackHandler, eventName string, args ...interface{}) {
	switch eventName {
	case "on_llm_start":
		// assuming the first arg is of type map[string]interface{}, and second is of type []string
		handler.OnLLMStart(args[0].(map[string]interface{}), args[1].([]string))
	// Add other cases as per your events
	default:
		return
	}
}
