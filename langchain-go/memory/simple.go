package memory

import "github.com/William-Bohm/langchain-go/langchain-go/memory/memorySchema"

type SimpleMemory struct {
	Memories map[string]interface{}
}

var _ memorySchema.BaseMemory = &SimpleMemory{}

func (s *SimpleMemory) MemoryVariables() ([]string, error) {
	keys := make([]string, 0, len(s.Memories))
	for k := range s.Memories {
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *SimpleMemory) LoadMemoryVariables(inputs map[string]interface{}) (map[string]interface{}, error) {
	return s.Memories, nil
}

func (s *SimpleMemory) SaveContext(inputs map[string]interface{}, outputs map[string]string) error {
	// Nothing should be saved or changed, my memory is set in stone.
	return nil
}

func (s *SimpleMemory) Clear() error {
	// Nothing to clear, got a memory like a vault.
	return nil
}
