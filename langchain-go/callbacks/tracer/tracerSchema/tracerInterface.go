package tracerSchema

type TracerInterface interface {
	AddChildRun(parentRun, childRun interface{}) error
	PersistRun(run interface{}) error
	PersistSession(sessionCreate TracerSession) (*TracerSession, error)
	GenerateId() *string
	NewSession(name string) (*TracerSession, error)
	LoadSession(sessionName string) (*TracerSession, error)
	LoadDefaultSession() (*TracerSession, error)
	Stack() []interface{}
	ExecutionOrder() int
	SetExecutionOrder(value int)
	Session() *TracerSession
	SetSession(value TracerSession)
}
