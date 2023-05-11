package chains

import (
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/memory/memorySchema"
	"path/filepath"
)

type Chain interface {
	ChainType() string
	InputKeys() []string
	OutputKeys() []string
	ValidateInputs(map[string]string) error
	ValidateOutputs(map[string]string) error
	Call(map[string]string) (map[string]string, error)
	Execute(map[string]interface{}, bool) (map[string]interface{}, error)
	PrepareOutputs(map[string]string, map[string]string, bool) map[string]string
	PrepareInputs(interface{}) (map[string]string, error)
	Apply([]map[string]interface{}) ([]map[string]string, error)
	Run(...interface{}) (string, error)
	ToDict() map[string]interface{}
	Save(string) error
}

type BaseChain struct {
	memory          memorySchema.BaseMemory
	callbackManager callbackSchema.BaseCallbackManager
	verbose         bool
}

func (bc *BaseChain) ChainType() string {
	panic("Saving not supported for this chain type.")
}

func (bc *BaseChain) ValidateInputs(inputs map[string]string) error {
	inputKeys := bc.InputKeys()
	missingKeys := []string{}

	for _, key := range inputKeys {
		if _, ok := inputs[key]; !ok {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("Missing some input keys: %v", missingKeys)
	}

	return nil
}

func (bc *BaseChain) ValidateOutputs(outputs map[string]string) error {
	outputKeys := bc.OutputKeys()
	if !equalStringSets(outputKeys, mapKeys(outputs)) {
		return fmt.Errorf("Did not get output keys that were expected. Got: %v. Expected: %v", mapKeys(outputs), outputKeys)
	}

	return nil
}
func (bc *BaseChain) Call(inputs map[string]string) (map[string]string, error) {
	panic("LLMChain logic not implemented.")
}

func (bc *BaseChain) Execute(
	inputs map[string]interface{},
	returnOnlyOutputs bool,
) (interface{}, error) {
	preparedInputs, err := bc.PrepareInputs(inputs)
	if err != nil {
		return nil, err
	}
	outputs, err := bc.Call(preparedInputs)
	if err != nil {
		return nil, err
	}
	return bc.PrepareOutputs(preparedInputs, outputs, returnOnlyOutputs), nil
}

func (bc *BaseChain) PrepareOutputs(
	inputs map[string]string,
	outputs map[string]string,
	returnOnlyOutputs bool,
) map[string]string {
	if returnOnlyOutputs {
		return outputs
	}
	merged := map[string]string{}
	for k, v := range inputs {
		merged[k] = v
	}
	for k, v := range outputs {
		merged[k] = v
	}
	return merged
}

func (bc *BaseChain) PrepareInputs(inputs interface{}) (map[string]string, error) {
	inputMap, ok := inputs.(map[string]string)
	if !ok {
		return nil, errors.New("inputs must be a map[string]string")
	}
	return inputMap, nil
}

func (bc *BaseChain) Apply(inputList []map[string]interface{}) ([]map[string]string, error) {
	outputList := []map[string]string{}
	for _, inputs := range inputList {
		output, err := bc.Execute(inputs, false)
		if err != nil {
			return nil, err
		}
		outputStr, ok := output.(map[string]string)
		if !ok {
			return nil, errors.New("output must be a map[string]string")
		}
		outputList = append(outputList, outputStr)
	}
	return outputList, nil
}

func (bc *BaseChain) ToDict() map[string]interface{} {
	return map[string]interface{}{
		"_type": bc.ChainType(),
	}
}

func (bc *BaseChain) Save(filePath string) error {
	ext := filepath.Ext(filePath)
	if ext != ".json" && ext != ".yaml" {
		return errors.New("file must be json or yaml")
	}
	return fmt.Errorf("Saving not supported for this chain type.")
}

func (bc *BaseChain) InputKeys() []string {
	panic("InputKeys method must be implemented by the LLMChain")
}

func (bc *BaseChain) OutputKeys() []string {
	panic("OutputKeys method must be implemented by the LLMChain")
}

func equalStringSets(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	set := map[string]struct{}{}
	for _, s := range a {
		set[s] = struct{}{}
	}
	for _, s := range b {
		if _, ok := set[s]; !ok {
			return false
		}
	}
	return true
}

func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
