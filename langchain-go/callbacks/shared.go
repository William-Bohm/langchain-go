package callbacks

import (
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"sync"
)

// SharedCallbackManager structure
type SharedCallbackManager struct {
	lock            sync.Mutex
	callbackManager callbackSchema.CallbackManager
	verbose         bool
}

// Global instance of SharedCallbackManager
var sharedCallbackManagerInstance *SharedCallbackManager
var once sync.Once

func GetSharedCallbackManagerInstance() *SharedCallbackManager {
	once.Do(func() {
		sharedCallbackManagerInstance = &SharedCallbackManager{
			callbackManager: callbackSchema.CallbackManager{},
		}
	})
	return sharedCallbackManagerInstance
}

func (s *SharedCallbackManager) OnLlmStart(serialized map[string]interface{}, prompts []string, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnLLMStart(serialized, prompts, s.verbose, kwargs)
}

func (s *SharedCallbackManager) OnLlmEnd(response llmSchema.LLMResult, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnLLMEnd(response, s.verbose, kwargs)
}

func (s *SharedCallbackManager) OnToolStart(serialized map[string]interface{}, inputStr string, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnToolStart(serialized, inputStr, s.verbose, kwargs)
}

func (s *SharedCallbackManager) OnToolEnd(output string, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnToolEnd(output, s.verbose, kwargs)
}

func (s *SharedCallbackManager) OnToolError(error error, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnToolError(error, s.verbose, kwargs)
}

func (s *SharedCallbackManager) AddHandler(callback callbackSchema.BaseCallbackHandler) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.AddHandler(callback)
}

func (s *SharedCallbackManager) RemoveHandler(callback callbackSchema.BaseCallbackHandler) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.RemoveHandler(callback)
}

func (s *SharedCallbackManager) SetHandlers(handlers []callbackSchema.BaseCallbackHandler) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.SetHandlers(handlers)
}

func (s *SharedCallbackManager) OnText(text string, verbose bool, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnText(text, verbose, kwargs)
}

func (s *SharedCallbackManager) OnAgentFinish(finish callbackSchema.AgentFinish, verbose bool, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnAgentFinish(finish, verbose, kwargs)
}

func (s *SharedCallbackManager) OnAgentAction(action callbackSchema.AgentAction, verbose bool, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnAgentAction(action, verbose, kwargs)
}

func (s *SharedCallbackManager) OnLlmError(error error, verbose bool, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnLLMError(error, verbose, kwargs)
}

func (s *SharedCallbackManager) OnChainStart(serialized map[string]interface{}, inputs map[string]interface{}, verbose bool, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnChainStart(serialized, inputs, verbose, kwargs)
}

func (s *SharedCallbackManager) OnChainEnd(outputs map[string]interface{}, verbose bool, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnChainEnd(outputs, verbose, kwargs)
}

func (s *SharedCallbackManager) OnChainError(error error, verbose bool, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnChainError(error, verbose, kwargs)
}

func (s *SharedCallbackManager) OnLlmNewToken(token string, verbose bool, kwargs map[string]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.callbackManager.OnLLMNewToken(token, verbose, kwargs)
}
