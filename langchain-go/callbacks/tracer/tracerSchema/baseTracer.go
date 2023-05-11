package tracerSchema

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"reflect"
	"time"
)

type BaseTracer struct {
	stack          []interface{}
	Session        *TracerSession
	executionOrder int
}

func NewBaseTracer(name string, executionOrder int, extra map[string]interface{}) *BaseTracer {
	session := NewTracerSession(name, extra)
	return &BaseTracer{
		executionOrder: executionOrder,
		Session:        session,
	}
}

func (bt *BaseTracer) AddChildRun(parentRun interface{}, childRun interface{}) error {
	// implementation goes here
	return nil
}

func (bt *BaseTracer) PersistRun(run interface{}) error {
	// implementation goes here
	return nil
}

func (bt *BaseTracer) PersistSession(sessionCreate TracerSession) (*TracerSession, error) {
	// implementation goes here
	return nil, nil
}

func (bt *BaseTracer) GenerateId() *string {
	// implementation goes here
	return nil
}

func (bt *BaseTracer) NewSession(name string) (*TracerSession, error) {
	// implementation goes here
	return nil, nil
}

func (bt *BaseTracer) LoadSession(sessionName string) (*TracerSession, error) {
	// implementation goes here
	return nil, nil
}

func (bt *BaseTracer) LoadDefaultSession() (*TracerSession, error) {
	// implementation goes here
	return nil, nil
}

func (bt *BaseTracer) GetStack() []interface{} {
	return bt.stack
}

func (bt *BaseTracer) ExecutionOrder() int {
	return bt.executionOrder
}

func (bt *BaseTracer) SetExecutionOrder(value int) {
	bt.executionOrder = value
}

func (bt *BaseTracer) GetSession() *TracerSession {
	return bt.Session
}

func (bt *BaseTracer) SetSession(value TracerSession) {
	bt.Session = &value
}

func (t *BaseTracer) StartTrace(run interface{}) error {
	t.executionOrder++

	if len(t.stack) > 0 {
		if _, ok := t.stack[len(t.stack)-1].(*callbackSchema.ChainRun); !ok {
			if _, ok := t.stack[len(t.stack)-1].(*callbackSchema.ToolRun); !ok {
				return errors.New("Nested " + reflect.TypeOf(run).Name() + " can only be logged inside a ChainRun or ToolRun")
			}
		}
		err := t.AddChildRun(t.stack[len(t.stack)-1], run)
		if err != nil {
			return err
		}
	}
	t.stack = append(t.stack, run)

	return nil
}

func (t *BaseTracer) EndTrace() {
	run := t.stack[len(t.stack)-1]
	t.stack = t.stack[:len(t.stack)-1]

	if len(t.stack) == 0 {
		t.executionOrder = 1
		err := t.PersistRun(run)
		if err != nil {
			return
		}
	}
}

func (t *BaseTracer) OnLLMStart(serialized map[string]interface{}, prompts []string, extra map[string]interface{}) error {
	if t.Session == nil {
		return errors.New("Initialize a Session with `new_session()` before starting a trace.")
	}

	llmRun := callbackSchema.NewLLMRun(
		t.executionOrder,
		serialized,
		t.Session.ID,
		prompts,
		&extra,
	)

	return t.StartTrace(llmRun)
}

func (t *BaseTracer) OnLLMNewToken(token string, extra map[string]interface{}) {
	// Handle a new token for an LLM run.
}

func (t *BaseTracer) OnLLMEnd(response *llmSchema.LLMResult, extra map[string]interface{}) error {
	if len(t.stack) == 0 || reflect.TypeOf(t.stack[len(t.stack)-1]).Name() != "LLMRun" {
		return errors.New("No LLMRun found to be traced")
	}

	llmRun := t.stack[len(t.stack)-1].(*callbackSchema.LLMRun)
	llmRun.EndTime = time.Now().UTC()
	llmRun.Response = response

	t.EndTrace()

	return nil
}

func (t *BaseTracer) OnLLMError(err error, extra map[string]interface{}) error {
	if len(t.stack) == 0 || reflect.TypeOf(t.stack[len(t.stack)-1]).Name() != "LLMRun" {
		return errors.New("No LLMRun found to be traced")
	}

	llmRun := t.stack[len(t.stack)-1].(*callbackSchema.LLMRun)
	llmRun.Error = &err
	llmRun.EndTime = time.Now().UTC()

	t.EndTrace()

	return nil
}

func (t *BaseTracer) OnChainStart(serialized map[string]interface{}, inputs map[string]interface{}, extra map[string]interface{}) error {
	if t.Session == nil { // ...
		return errors.New("Initialize a Session with `new_session()` before starting a trace.")
	}

	chainRun := callbackSchema.NewChainRun(
		t.executionOrder,
		serialized,
		t.Session.ID,
		inputs,
		&extra,
	)

	return t.StartTrace(chainRun)
}

func (t *BaseTracer) OnChainEnd(outputs map[string]interface{}, extra map[string]interface{}) error {
	if len(t.stack) == 0 || reflect.TypeOf(t.stack[len(t.stack)-1]).Name() != "ChainRun" {
		return errors.New("No ChainRun found to be traced")
	}

	chainRun := t.stack[len(t.stack)-1].(*callbackSchema.ChainRun)
	chainRun.EndTime = time.Now().UTC()
	chainRun.Outputs = &outputs

	t.EndTrace()

	return nil
}

func (t *BaseTracer) OnChainError(err error, extra map[string]interface{}) error {
	if len(t.stack) == 0 || reflect.TypeOf(t.stack[len(t.stack)-1]).Name() != "ChainRun" {
		return errors.New("No ChainRun found to be traced")
	}

	chainRun := t.stack[len(t.stack)-1].(*callbackSchema.ChainRun)
	chainRun.Error = &err
	chainRun.EndTime = time.Now().UTC()

	t.EndTrace()

	return nil
}

func (t *BaseTracer) OnToolStart(serialized map[string]interface{}, inputStr string, extra map[string]interface{}) error {
	if t.Session == nil {
		return errors.New("Initialize a Session with `new_session()` before starting a trace.")
	}

	toolRun := callbackSchema.NewToolRun(
		t.executionOrder,
		serialized,
		t.Session.ID,
		inputStr,
		&extra,
	)

	return t.StartTrace(toolRun)
}

func (t *BaseTracer) OnToolEnd(output string, extra map[string]interface{}) error {
	if len(t.stack) == 0 || reflect.TypeOf(t.stack[len(t.stack)-1]).Name() != "ToolRun" {
		return errors.New("No ToolRun found to be traced")
	}

	toolRun := t.stack[len(t.stack)-1].(*callbackSchema.ToolRun)
	toolRun.EndTime = time.Now().UTC()
	toolRun.Output = &output

	t.EndTrace()

	return nil
}

func (t *BaseTracer) OnToolError(err error, extra map[string]interface{}) error {
	if len(t.stack) == 0 || reflect.TypeOf(t.stack[len(t.stack)-1]).Name() != "ToolRun" {
		return errors.New("No ToolRun found to be traced")
	}

	toolRun := t.stack[len(t.stack)-1].(*callbackSchema.ToolRun)
	toolRun.Error = &err
	toolRun.EndTime = time.Now().UTC()

	t.EndTrace()

	return nil
}
