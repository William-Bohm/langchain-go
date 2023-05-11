package utils

import (
	"fmt"
)

func GetPromptInputKey(inputs map[string]interface{}, memoryVariables []string) (string, error) {
	stop := "stop"

	promptInputKeys := make(map[string]struct{})
	for k := range inputs {
		promptInputKeys[k] = struct{}{}
	}

	for _, mv := range memoryVariables {
		delete(promptInputKeys, mv)
	}
	delete(promptInputKeys, stop)

	if len(promptInputKeys) != 1 {
		return "", fmt.Errorf("one input key expected got %v", promptInputKeys)
	}

	var key string
	for k := range promptInputKeys {
		key = k
		break
	}
	return key, nil
}
