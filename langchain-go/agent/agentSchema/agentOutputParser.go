package agentSchema

type AgentOutputParser interface {
	Parse(text string) interface{}
}
