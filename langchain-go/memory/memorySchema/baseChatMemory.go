package memorySchema

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/memory/chatMessageHistories"
	"github.com/William-Bohm/langchain-go/langchain-go/memory/utils"
)

type BaseChatMemoryInterface interface {
	BaseMemory
	GetInputOutput(inputs map[string]interface{}, outputs map[string]string) (string, string, error)
}

type BaseChatMemory struct {
	BaseMemory
	ChatMemory     chatMessageHistories.BaseChatMessageHistory
	OutputKey      *string
	InputKey       *string
	ReturnMessages bool
}

// GetInputOutput retrieve the input and output messages in the conversatoin history
func (bcm *BaseChatMemory) GetInputOutput(inputs map[string]interface{}, outputs map[string]string) (string, string, error) {
	var promptInputKey string
	if bcm.InputKey == nil {
		memoryVariables, err := bcm.MemoryVariables()
		promptInputKey, err = utils.GetPromptInputKey(inputs, memoryVariables)
		if err != nil {
			return "", "", err
		}
	} else {
		promptInputKey = *bcm.InputKey
	}
	var outputKey string
	if bcm.OutputKey == nil {
		if len(outputs) != 1 {
			return "", "", fmt.Errorf("One output key expected, got %v", outputs)
		}
		for k := range outputs {
			outputKey = k
			break
		}
	} else {
		outputKey = *bcm.OutputKey
	}
	return inputs[promptInputKey].(string), outputs[outputKey], nil
}

func (bcm *BaseChatMemory) SaveContext(inputs map[string]interface{}, outputs map[string]string) error {
	inputStr, outputStr, err := bcm.GetInputOutput(inputs, outputs)
	if err != nil {
		return err
	}
	bcm.ChatMemory.AddUserMessage(inputStr)
	bcm.ChatMemory.AddAIMessage(outputStr)
	return nil
}

func (bcm *BaseChatMemory) Clear() error {
	return bcm.ChatMemory.Clear()
}
