package schema

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
	"sync"
)

var langchainVerbose bool

func getLangchainVerbose() bool {
	return langchainVerbose
}

type BaseChatModel struct {
	verbose               bool
	callbackManager       callbackSchema.BaseCallbackManager
	arbitraryTypesAllowed bool
	extra                 string
}

func NewBaseChatModel(verbose bool, cm callbackSchema.BaseCallbackManager) *BaseChatModel {
	if cm == nil {
		cm = callbackSchema.GetCallbackManager()
	}
	return &BaseChatModel{
		verbose:               verbose,
		callbackManager:       cm,
		arbitraryTypesAllowed: true,
		extra:                 "forbid",
	}
}

func (m *BaseChatModel) combineLLMOutputs(llmOutputs []map[string]interface{}) map[string]interface{} {
	return make(map[string]interface{})
}

func (m *BaseChatModel) Generate(messages [][]rootSchema.BaseMessage, stop []string) (llmSchema.LLMResult, error) {
	var wg sync.WaitGroup
	wg.Add(len(messages))

	results := make([]schema.ChatResult, len(messages))
	for i, msg := range messages {
		go func(i int, msg []rootSchema.BaseMessage) {
			defer wg.Done()
			results[i] = m._generate(msg, stop)
		}(i, msg)
	}
	wg.Wait()

	llmOutput := m.combineLLMOutputs([]map[string]interface{}{})
	var generations []schema.ChatGeneration
	for _, res := range results {
		generations = append(generations, res.generations...)
	}

	return llmSchema.LLMResult{
		generations: generations,
		llmOutput:   llmOutput,
	}, nil
}

func (m *BaseChatModel) GeneratePrompt(prompts []schema.PromptValue, stop []string) (llmSchema.LLMResult, error) {
	var promptMessages [][]rootSchema.BaseMessage
	var promptStrings []string
	for _, p := range prompts {
		promptMessages = append(promptMessages, p.ToMessages())
		promptStrings = append(promptStrings, p.ToString())
	}
	m.callbackManager.OnLLMStart(map[string]interface{}{"name": "BaseChatModel"}, promptStrings, m.verbose)

	defer func() {
		if err := recover(); err != nil {
			m.callbackManager.OnLLMError(err.(error), m.verbose)
			panic(err)
		}
	}()

	output, err := m.Generate(promptMessages, stop)
	if err != nil {
		m.callbackManager.OnLLMError(err, m.verbose)
		return schema.LLMResult{}, err
	}

	m.callbackManager.OnLLMEnd(output, m.verbose)
	return output, nil
}

func (m *BaseChatModel) _generate(messages []rootSchema.BaseMessage, stop []string) schema.ChatResult {
	panic(errors.New("_generate not implemented"))
}

func (m *BaseChatModel) Call(messages []rootSchema.BaseMessage, stop []string) rootSchema.BaseMessage {
	return m._generate(messages, stop).generations[0].message
}

type SimpleChatModel struct {
	BaseChatModel
}

func NewSimpleChatModel(verbose bool, cm callbackSchema.BaseCallbackManager) *SimpleChatModel {
	base := NewBaseChatModel(verbose, cm)
	return &SimpleChatModel{*base}
}

func (m *SimpleChatModel) _generate(messages []rootSchema.BaseMessage, stop []string) schema.ChatResult {
	outputStr := m._call(messages, stop)
	message := schema.AIMessage{content: outputStr}
	generation := schema.ChatGeneration{message: message}

	return schema.ChatResult{generations: []schema.ChatGeneration{generation}}
}

func (m *SimpleChatModel) _call(messages []schema.BaseMessage, stop []string) string {
	panic(errors.New("_call not implemented"))
}
