package memory

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/chains"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
)

type SummarizerMixin struct {
	HumanPrefix       string
	AIPrefix          string
	LLM               llmSchema.BaseLanguageModel
	Prompt            promptSchema.BasePromptTemplate
	SummaryMessageCls rootSchema.BaseMessage
}

func (s *SummarizerMixin) PredictNewSummary(messages []rootSchema.BaseMessageInterface, existingSummary string) (string, error) {
	newLines, err := rootSchema.GetBufferString(messages, s.HumanPrefix, s.AIPrefix)
	if err != nil {
		return "", err
	}

	chain := chains.LLMChain{
		LLM:    s.LLM,
		Prompt: s.Prompt,
	}
	return chain.Predict(existingSummary, newLines)
}

type ConversationSummaryMemory struct {
	BaseChatMemory
	SummarizerMixin
	Buffer    string
	MemoryKey string
}

func (c *ConversationSummaryMemory) MemoryVariables() ([]string, error) {
	return []string{c.MemoryKey}, nil
}

func (c *ConversationSummaryMemory) LoadMemoryVariables(inputs map[string]interface{}) (map[string]interface{}, error) {
	var buffer interface{}
	if c.ReturnMessages {
		buffer = []llmSchema.BaseMessage{c.SummaryMessageCls}
	} else {
		buffer = c.Buffer
	}
	return map[string]interface{}{
		c.MemoryKey: buffer,
	}, nil
}

func ValidatePromptInputVariables(prompt promptSchema.BasePromptTemplate) error {
	promptVariables := prompt.InputVariables()
	expectedKeys := []string{"summary", "new_lines"}
	if !equalStringSlices(promptVariables, expectedKeys) {
		return fmt.Errorf("got unexpected prompt input variables. The prompt expects %v, but it should have %v", promptVariables, expectedKeys)
	}
	return nil
}

func (c *ConversationSummaryMemory) SaveContext(inputs map[string]interface{}, outputs map[string]string) error {
	err := c.BaseChatMemory.SaveContext(inputs, outputs)
	if err != nil {
		return err
	}
	newSummary, err := c.PredictNewSummary(c.ChatMemory.Messages()[-2:], c.Buffer)
	if err != nil {
		return err
	}
	c.Buffer = newSummary
	return nil
}

func (c *ConversationSummaryMemory) Clear() error {
	err := c.BaseChatMemory.Clear()
	if err != nil {
		return err
	}
	c.Buffer = ""
	return nil
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
