package tracerSchema

import (
	"errors"
	"sync"
)

type Tracer struct {
	*BaseTracer
	mutex sync.Mutex
}

func NewTracer(name string, executionOrder int, extra map[string]interface{}) *Tracer {
	baseTracer := NewBaseTracer(name, executionOrder, extra)
	return &Tracer{
		BaseTracer: baseTracer,
	}
}

func (t *Tracer) getStack() []interface{} {
	return t.stack
}

func (t *Tracer) getExecutionOrder() int {
	return t.executionOrder
}

func (t *Tracer) setExecutionOrder(value int) {
	t.executionOrder = value
}

func (t *Tracer) getSession() *TracerSession {
	return t.Session
}

func (t *Tracer) setSession(value *TracerSession) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if len(t.stack) > 0 {
		return errors.New("cannot set a Session while a trace is being recorded")
	}

	t.Session = value
	return nil
}
