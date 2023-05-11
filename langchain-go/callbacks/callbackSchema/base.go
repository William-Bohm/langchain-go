package callbackSchema

import "github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"

type BaseCallbackHandler interface {
	AlwaysVerbose() bool
	IgnoreLLM() bool
	IgnoreChain() bool
	IgnoreAgent() bool
	OnLLMStart(serialized map[string]interface{}, prompts []string, kwargs map[string]interface{}) (interface{}, error)
	OnLLMNewToken(token string, kwargs map[string]interface{}) (interface{}, error)
	OnLLMEnd(response llmSchema.LLMResult, kwargs map[string]interface{}) (interface{}, error)
	OnLLMError(err error, kwargs map[string]interface{}) (interface{}, error)
	OnChainStart(serialized map[string]interface{}, inputs map[string]interface{}, kwargs map[string]interface{}) (interface{}, error)
	OnChainEnd(outputs map[string]interface{}, kwargs map[string]interface{}) (interface{}, error)
	OnChainError(err error, kwargs map[string]interface{}) (interface{}, error)
	OnToolStart(serialized map[string]interface{}, inputStr string, kwargs map[string]interface{}) (interface{}, error)
	OnToolEnd(output string, kwargs map[string]interface{}) (interface{}, error)
	OnToolError(err error, kwargs map[string]interface{}) (interface{}, error)
	OnText(text string, kwargs map[string]interface{}) (interface{}, error)
	OnAgentAction(action AgentAction, kwargs map[string]interface{}) (interface{}, error)
	OnAgentFinish(finish AgentFinish, kwargs map[string]interface{}) (interface{}, error)
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

func (c *CallbackManager) OnLLMStart(serialized map[string]interface{}, prompts []string, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreLLM() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnLLMStart(serialized, prompts, kwargs)
			}
		}
	}
}

func (c *CallbackManager) OnLLMNewToken(token string, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreLLM() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnLLMNewToken(token, kwargs)
			}
		}
	}
}

func (c *CallbackManager) OnLLMEnd(response llmSchema.LLMResult, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreLLM() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnLLMEnd(response, kwargs)
			}
		}
	}
}

func (c *CallbackManager) OnLLMError(err error, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreLLM() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnLLMError(err, kwargs)
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

func (c *CallbackManager) OnChainEnd(outputs map[string]interface{}, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreChain() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnChainEnd(outputs, kwargs)
			}
		}
	}
}

func (c *CallbackManager) OnChainError(err error, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreChain() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnChainError(err, kwargs)
			}
		}
	}
}

func (c *CallbackManager) OnToolStart(serialized map[string]interface{}, input_str string, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreAgent() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnToolStart(serialized, input_str, kwargs)
			}
		}
	}
}

func (c *CallbackManager) OnAgentAction(action AgentAction, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreAgent() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnAgentAction(action, kwargs)
			}
		}
	}
}

func (c *CallbackManager) OnToolEnd(output string, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreAgent() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnToolEnd(output, kwargs)
			}
		}
	}
}

func (c *CallbackManager) OnToolError(err error, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreAgent() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnToolError(err, kwargs)
			}
		}
	}
}

func (c *CallbackManager) OnText(text string, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if verbose || handler.AlwaysVerbose() {
			handler.OnText(text, kwargs)
		}
	}
}

func (c *CallbackManager) OnAgentFinish(finish AgentFinish, verbose bool, kwargs map[string]interface{}) {
	for _, handler := range c.handlers {
		if !handler.IgnoreAgent() {
			if verbose || handler.AlwaysVerbose() {
				handler.OnAgentFinish(finish, kwargs)
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
func handleEvent(handler BaseCallbackHandler, eventName string, kwargs map[string]interface{}) {
	switch eventName {
	case "on_llm_start":
		// assuming the first arg is of type map[string]interface{}, and second is of type []string
		handler.OnLLMStart(args[0].(map[string]interface{}), args[1].([]string))
	// Add other cases as per your events
	default:
		return
	}
}
