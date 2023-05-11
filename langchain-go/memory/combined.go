package memory

import (
	"github.com/William-Bohm/langchain-go/langchain-go/memory/memorySchema"
)

type CombinedMemory struct {
	Memories []memorySchema.BaseMemory
}

func NewCombinedMemory(memories []memorySchema.BaseMemory) *CombinedMemory {
	return &CombinedMemory{
		Memories: memories,
	}
}

func (c *CombinedMemory) MemoryVariables() ([]string, error) {
	memoryVariables := []string{}

	for _, memory := range c.Memories {
		newMemoryVariables, err := memory.MemoryVariables()
		if err != nil {
			return nil, err
		}
		memoryVariables = append(memoryVariables, newMemoryVariables...)
	}

	return memoryVariables, nil
}

func (c *CombinedMemory) LoadMemoryVariables(inputs map[string]interface{}) (map[string]interface{}, error) {
	memoryData := make(map[string]interface{})

	for _, memory := range c.Memories {
		data, err := memory.LoadMemoryVariables(inputs)
		if err != nil {
			return nil, err
		}
		for k, v := range data {
			memoryData[k] = v
		}
	}

	return memoryData, nil
}

func (c *CombinedMemory) SaveContext(inputs map[string]interface{}, outputs map[string]string) {
	for _, memory := range c.Memories {
		memory.SaveContext(inputs, outputs)
	}
}

func (c *CombinedMemory) Clear() {
	for _, memory := range c.Memories {
		memory.Clear()
	}
}
