package callbackSchema

type AgentAction struct {
	Tool      string
	ToolInput interface{}
	Log       string
}

type AgentFinish struct {
	ReturnValues map[string]interface{}
	Log          string
}
