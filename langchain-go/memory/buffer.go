package memory

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/memory/chatMessageHistories"
	"github.com/William-Bohm/langchain-go/langchain-go/memory/memorySchema"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
)

type ConversationBufferMemory struct {
	memorySchema.BaseMemory
	HumanPrefix string
	AIPrefix    string
	MemoryKey   string
	ChatMemory  chatMessageHistories.ChatMessageHistory
}

func NewConversationBufferMemory(chatMemory *chatMessageHistories.ChatMessageHistory) *ConversationBufferMemory {
	return &ConversationBufferMemory{
		HumanPrefix: "Human",
		AIPrefix:    "AI",
		MemoryKey:   "history",
		ChatMemory:  *chatMemory,
	}
}

func (c *ConversationBufferMemory) Buffer() string {
	messages, _ := c.ChatMemory.Messages()
	return getBufferString(messages, c.HumanPrefix, c.AIPrefix)
}

func (c *ConversationBufferMemory) MemoryVariables() ([]string, error) {
	return []string{c.MemoryKey}, nil
}

func (c *ConversationBufferMemory) LoadMemoryVariables(inputs map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{c.MemoryKey: c.Buffer()}, nil
}

type ConversationStringBufferMemory struct {
	HumanPrefix string
	AIPrefix    string
	Buffer      string
	OutputKey   *string
	InputKey    *string
	MemoryKey   string
}

func NewConversationStringBufferMemory() *ConversationStringBufferMemory {
	return &ConversationStringBufferMemory{
		HumanPrefix: "Human",
		AIPrefix:    "AI",
		MemoryKey:   "history",
	}
}

func (c *ConversationStringBufferMemory) MemoryVariables() []string {
	return []string{c.MemoryKey}
}

func (c *ConversationStringBufferMemory) LoadMemoryVariables(inputs map[string]interface{}) (map[string]string, error) {
	return map[string]string{c.MemoryKey: c.Buffer}, nil
}

func (c *ConversationStringBufferMemory) SaveContext(inputs map[string]interface{}, outputs map[string]string) {
	promptInputKey := c.InputKey
	if promptInputKey == nil {
		promptInputKey = getPromptInputKey(inputs, c.MemoryVariables())
	}
	outputKey := c.OutputKey
	if outputKey == nil {
		if len(outputs) != 1 {
			panic(fmt.Sprintf("One output key expected, got %v", outputs))
		}
		for k := range outputs {
			outputKey = &k
			break
		}
	}

	human := fmt.Sprintf("%s: %s", c.HumanPrefix, inputs[*promptInputKey].(string))
	ai := fmt.Sprintf("%s: %s", c.AIPrefix, outputs[*outputKey])

	c.Buffer += "\n" + human + "\n" + ai
}

func (c *ConversationStringBufferMemory) Clear() {
	c.Buffer = ""
}

func getBufferString(messages []rootSchema.BaseMessageInterface, humanPrefix string, aiPrefix string) string {
	buffer := ""
	for _, message := range messages {
		prefix := ""
		switch message.(type) {
		case *rootSchema.HumanMessage:
			prefix = humanPrefix
		case *rootSchema.AIMessage:
			prefix = aiPrefix
		}
		if buffer != "" {
			buffer += "\n"
		}
		buffer += fmt.Sprintf("%s: %s", prefix, message.GetContent())
	}
	return buffer
}

func getPromptInputKey(inputs map[string]interface{}, memoryVariables []string) *string {
	for k := range inputs {
		found := false
		for _, v := range memoryVariables {
			if k == v {
				found = true
				break
			}
		}
		if !found {
			return &k
		}
	}
	return nil
}
