package memory

import (
	"github.com/William-Bohm/langchain-go/langchain-go/memory/memorySchema"
)

type ReadOnlySharedMemory struct {
	Memory memorySchema.BaseMemory
}

func (r *ReadOnlySharedMemory) MemoryVariables() ([]string, error) {
	return r.Memory.MemoryVariables()
}

func (r *ReadOnlySharedMemory) LoadMemoryVariables(inputs map[string]interface{}) (map[string]interface{}, error) {
	return r.Memory.LoadMemoryVariables(inputs)
}

func (r *ReadOnlySharedMemory) SaveContext(inputs map[string]interface{}, outputs map[string]string) error {
	// Nothing should be saved or changed
	return nil
}

func (r *ReadOnlySharedMemory) Clear() error {
	// Nothing to clear, got a memory like a vault.
	return nil
}
