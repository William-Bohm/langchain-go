package memory

import (
	"github.com/William-Bohm/langchain-go/langchain-go/memory/memorySchema"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
)

type ConversationBufferWindowMemory struct {
	memorySchema.BaseChatMemory
	HumanPrefix string
	AIPrefix    string
	MemoryKey   string
	K           int
}

func NewConversationBufferWindowMemory(chatMemory memorySchema.BaseChatMemory, humanPrefix, aiPrefix, memoryKey string, k int) *ConversationBufferWindowMemory {
	return &ConversationBufferWindowMemory{
		BaseChatMemory: chatMemory,
		HumanPrefix:    humanPrefix,
		AIPrefix:       aiPrefix,
		MemoryKey:      memoryKey,
		K:              k,
	}
}

func (c *ConversationBufferWindowMemory) Buffer() ([]rootSchema.BaseMessageInterface, error) {
	messages, err := c.ChatMemory.Messages()
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *ConversationBufferWindowMemory) MemoryVariables() []string {
	return []string{c.MemoryKey}
}

func (c *ConversationBufferWindowMemory) LoadMemoryVariables(inputs map[string]interface{}) (map[string]string, error) {
	buffer, err := c.Buffer()
	if err != nil {
		return nil, err
	}
	if c.K > 0 {
		buffer = buffer[-c.K*2:]
	} else {
		buffer = []rootSchema.BaseMessageInterface{}
	}

	var result string
	if !c.ReturnMessages {
		result = getBufferString(buffer, c.HumanPrefix, c.AIPrefix)
	}

	return map[string]string{c.MemoryKey: result}, nil
}
