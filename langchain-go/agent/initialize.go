package agent

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/agent/agentSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/tools/toolSchema"
)

func InitializeAgent(
	tools []toolSchema.BaseTool,
	llm llmSchema.BaseLanguageModel,
	agent AgentType,
	callbackManager callbackSchema.BaseCallbackManager,
	agentPath string,
	agentKwargs map[string]interface{},
	kwargs map[string]interface{},
) (*agentSchema.AgentExecutor, error) {
	if agent == "" && agentPath == "" {
		agent = ZERO_SHOT_REACT_DESCRIPTION
	}
	if agent != "" && agentPath != "" {
		return nil, errors.New("both `agent` and `agent_path` are specified, but at most only one should be")
	}
	var agentObj agentSchema.BaseAgent
	if agent != "" {
		agentClass, ok := AGENT_TO_CLASS[agent]
		if !ok {
			return nil, errors.New(string("unknown agent type: " + agent))
		}
		if agentKwargs == nil {
			agentKwargs = make(map[string]interface{})
		}
		var err error
		agentObj, err = agentClass.FromLLMAndTools(llm, tools, callbackManager, agentKwargs)
		if err != nil {
			return nil, err
		}
	} else if agentPath != "" {
		var err error
		agentObj, err = LoadAgent(agentPath, kwargs)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("somehow both `agent` and `agent_path` are None, this should never happen")
	}
	return agentSchema.NewAgentExecutor(agentObj, tools, callbackManager, false), nil
}
