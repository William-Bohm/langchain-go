package tracerSchema

import (
	"github.com/google/uuid"
	"time"
)

type TracerSession struct {
	StartTime time.Time
	Name      string
	Extra     map[string]interface{}
	ID        uuid.UUID
}

func NewTracerSession(name string, extra map[string]interface{}) *TracerSession {
	return &TracerSession{
		Name:  name,
		Extra: extra,
		ID:    uuid.New(),
	}
}
