package toolSchema

import (
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"reflect"
	"strings"
)

type BaseToolInterface interface {
	Args() map[string]interface{}
	ParseInput(toolInput interface{}) error
	Run(toolInput interface{}, verbose *bool, startColor, color string, kwargs map[string]interface{}) (string, error)
	Call(toolInput string) (string, error)
	// abstract
	run(args ...interface{}) (string, error)
}

type BaseTool struct {
	BaseToolInterface
	Name            string
	Description     string
	ArgsSchema      map[string]interface{}
	ReturnDirect    bool
	Verbose         bool
	CallbackManager callbackSchema.BaseCallbackManager
}

func (b *BaseTool) Args() map[string]interface{} {
	if b.ArgsSchema != nil {
		return b.ArgsSchema
	} else {
		return map[string]interface{}{}
	}
}

func (b *BaseTool) ParseInput(toolInput interface{}) error {
	inputArgs := b.ArgsSchema
	switch toolInput := toolInput.(type) {
	case string:
		if inputArgs != nil {
			key_ := reflect.TypeOf(inputArgs).Elem().Field(0).Name
			inputArgs[key_] = toolInput
		}
	case map[string]interface{}:
		if inputArgs != nil {
			inputArgs = toolInput
		}
	default:
		return errors.New("toolInput must be a string or map[string]interface{}")
	}
	return nil
}

func (b *BaseTool) Run(toolInput interface{}, verbose *bool, color string, kwargs map[string]interface{}) (string, error) {
	err := b.ParseInput(toolInput)
	if err != nil {
		return "", err
	}
	verbose_ := b.Verbose
	if !b.Verbose && verbose != nil {
		verbose_ = *verbose
	}
	toolStartCallback := make(map[string]interface{})
	toolStartCallback["verbose"] = verbose_
	toolStartCallback["color"] = color
	toolStartCallback["args"] = kwargs
	b.CallbackManager.OnToolStart(
		map[string]interface{}{"name": b.Name, "description": b.Description},
		strings.Trim(strings.Replace(fmt.Sprint(toolInput), " ", ",", -1), "[]"),
		toolStartCallback,
	)
	var observation string
	var toolArgs []interface{}
	var toolKwargs map[string]interface{}
	switch toolInput := toolInput.(type) {
	case string:
		toolArgs = append(toolArgs, toolInput)
	case map[string]interface{}:
		toolKwargs = toolInput
	default:
		return "", errors.New("toolInput must be a string or map[string]interface{}")
	}
	observation, err = b.run(toolArgs, toolKwargs)
	toolErrorCallback := make(map[string]interface{})
	toolErrorCallback["verbose"] = verbose_
	if err != nil {
		b.CallbackManager.OnToolError(err, toolErrorCallback)
		return "", err
	}
	callback := make(map[string]interface{})
	callback["color"] = color
	callback["verbose"] = verbose_
	callback["name"] = b.Name
	callback["args"] = kwargs
	b.CallbackManager.OnToolEnd(observation, callback)
	return observation, nil
}

func (b *BaseTool) Call(toolInput string) (string, error) {
	return b.Run(toolInput, nil, "green", "green", map[string]interface{}{})
}

func toArgsAndKwargs(runInput interface{}) ([]interface{}, map[string]interface{}) {
	switch runInput := runInput.(type) {
	case string:
		return []interface{}{runInput}, map[string]interface{}{}
	case map[string]interface{}:
		return []interface{}{}, runInput
	default:
		return []interface{}{}, map[string]interface{}{}
	}
}

//func main() {
//	// Example usage
//	tool := &BaseTool{
//		Name:            "Tool 1",
//		Description:     "This is tool 1",
//		ArgsSchema:      map[string]interface{}{},
//		ReturnDirect:    false,
//		Verbose:         false,
//		CallbackManager: getCallbackManager(),
//	}
//
//	tool.Run("input", nil, "green", "green", map[string]interface{}{})
//}
