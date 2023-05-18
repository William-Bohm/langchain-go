package agentSchema

type AgentAction struct {
	Tool      string
	ToolInput interface{}
	Log       string
}

type AgentFinish struct {
	ReturnValues map[string]interface{}
	Log          string
}

type AgentStep struct {
	AgentAction
	Observation string
}

type IntermediateStep struct {
	AgentAction
	Output string
}
