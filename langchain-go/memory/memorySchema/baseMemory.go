package memorySchema

// BaseMemory Base interface for memory in chains.
type BaseMemory interface {
	// MemoryVariables Input keys this memory class will load dynamically.
	MemoryVariables() ([]string, error)
	// LoadMemoryVariables Return key-value pairs given the text input to the chain.
	LoadMemoryVariables(inputs map[string]interface{}) (map[string]interface{}, error)
	// SaveContext Save the context of this model run to memory.
	SaveContext(inputs map[string]interface{}, outputs map[string]string) error
	// Clear memory contents.
	Clear() error
}
