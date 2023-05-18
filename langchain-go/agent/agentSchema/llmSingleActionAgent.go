package agentSchema

import "github.com/William-Bohm/langchain-go/langchain-go/chains"

type LLMSingleActionAgent struct {
	BaseSingleActionAgent
	LLMChain     chains.LLMChain
	OutputParser AgentOutputParser
	Stop         []string
}

func (agent *LLMSingleActionAgent) InputKeys() []string {
	inputKeys := agent.LLMChain.InputKeys()
	for i := len(inputKeys) - 1; i >= 0; i-- {
		if inputKeys[i] == "intermediate_steps" {
			inputKeys = append(inputKeys[:i], inputKeys[i+1:]...)
		}
	}
	return inputKeys
}

func (agent *LLMSingleActionAgent) Plan(intermediateSteps []struct {
	AgentAction
	string
}, kwargs map[string]interface{}) (AgentAction, AgentFinish) {
	output, _ := agent.LLMChain.Run(intermediateSteps, agent.Stop, kwargs)
	return agent.OutputParser.Parse(output)
}

func (agent *LLMSingleActionAgent) ToolRunLoggingKwargs() map[string]interface{} {
	kwargs := make(map[string]interface{})
	kwargs["llm_prefix"] = ""
	if len(agent.Stop) == 0 {
		kwargs["observation_prefix"] = ""
	} else {
		kwargs["observation_prefix"] = agent.Stop[0]
	}
	return kwargs
}
