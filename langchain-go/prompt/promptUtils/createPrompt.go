package promptUtils

import (
	"fmt"
	"strings"
)

// AddInputVariablesToPrompt replaces {variables} with desired values in a string
func AddInputVariablesToPrompt(inputVariables map[string]interface{}, template string) string {
	for key, value := range inputVariables {
		variable := fmt.Sprintf("{%s}", key)
		strValue := fmt.Sprintf("%v", value)
		template = strings.ReplaceAll(template, variable, strValue)
	}
	return template
}
