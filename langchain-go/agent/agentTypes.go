package agent

type AgentType string

const (
	zeroShotReactDescription           AgentType = "zero-shot-react-description"
	reactDocstore                      AgentType = "react-docstore"
	selfAskWithSearch                  AgentType = "self-ask-with-search"
	conversationalReactDescription     AgentType = "conversational-react-description"
	chatZeroShotReactDescription       AgentType = "chat-zero-shot-react-description"
	chatConversationalReactDescription AgentType = "chat-conversational-react-description"
)
