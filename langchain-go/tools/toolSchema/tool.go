package toolSchema

import (
	"errors"
)

type CallableFunc func(args ...interface{}) string
type CallableCoroutine func(args ...interface{}) (string, error)

type Tool struct {
	BaseTool
	Description string
	Func        CallableFunc
	Coroutine   CallableCoroutine
}

func NewTool(name string, fn CallableFunc, description string, kwargs ...interface{}) *Tool {
	base := &BaseTool{Name: name}
	return &Tool{BaseTool: *base, Description: description, Func: fn}
}

func (t *Tool) Args() map[string]interface{} {
	if t.ArgsSchema != nil {
		return t.ArgsSchema["properties"].(map[string]interface{})
	} else {
		return nil
	}
}

func (t *Tool) Run(args ...interface{}) string {
	return t.Func(args...)
}

func (t *Tool) ARun(args ...interface{}) (string, error) {
	if t.Coroutine != nil {
		return t.Coroutine(args...)
	} else {
		return "", errors.New("Tool does not support async")
	}
}
