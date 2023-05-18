package agentSchema

import (
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/chains"
	"github.com/William-Bohm/langchain-go/langchain-go/memory"
	"github.com/William-Bohm/langchain-go/langchain-go/memory/chatMessageHistories"
	"github.com/William-Bohm/langchain-go/langchain-go/tools"
	"github.com/William-Bohm/langchain-go/langchain-go/tools/toolSchema"
	"time"
)

type AgentExecutor struct {
	chains.Chain
	chains.BaseChain
	agent                   BaseAgent
	tools                   []toolSchema.BaseTool
	returnIntermediateSteps bool
	maxIterations           int
	maxExecutionTime        float64
	earlyStoppingMethod     string
}

func NewAgentExecutor(agent BaseAgent, tools []toolSchema.BaseTool, callbackManager callbackSchema.BaseCallbackManager, verbose bool) *AgentExecutor {
	return &AgentExecutor{
		BaseChain: chains.BaseChain{
			Memory:          memory.NewConversationBufferMemory(chatMessageHistories.NewChatMessageHistory()),
			CallbackManager: callbackManager,
			Verbose:         verbose,
		},
		agent:                   agent,
		tools:                   tools,
		returnIntermediateSteps: false,
		maxIterations:           15,
		maxExecutionTime:        0,
		earlyStoppingMethod:     "force",
	}
}

func (a *AgentExecutor) Save(filePath string) error {
	return errors.New("Saving not supported for agent executors. If you are trying to save the agent, please use the `.save_agent(...)`")
}

func (a *AgentExecutor) SaveAgent(filePath string) {
	a.agent.(BaseAgent).Save(filePath)
}

func (a *AgentExecutor) InputKeys() []string {
	return a.agent.(BaseAgent).InputKeys()
}

func (a *AgentExecutor) OutputKeys() []string {
	if a.returnIntermediateSteps {
		return append(a.agent.(BaseAgent).ReturnValues(), "intermediate_steps")
	}
	return a.agent.(BaseAgent).ReturnValues()
}

func (a *AgentExecutor) LookupTool(name string) (toolSchema.BaseTool, error) {
	for _, tool := range a.tools {
		if tool.Name == name {
			return tool, nil
		}
	}
	return toolSchema.BaseTool{}, fmt.Errorf("tool %s not found", name)
}

func (a *AgentExecutor) ShouldContinue(iterations int, timeElapsed float64) bool {
	if a.maxIterations != 0 && iterations >= a.maxIterations {
		return false
	}
	if a.maxExecutionTime != 0 && timeElapsed >= a.maxExecutionTime {
		return false
	}
	return true
}

func (a *AgentExecutor) TakeNextStep(nameToToolMap map[string]toolSchema.BaseTool, colorMapping map[string]string, inputs map[string]interface{}, intermediateSteps []IntermediateStep) (interface{}, error) {
	output := a.agent.(BaseAgent).Plan(intermediateSteps, inputs)
	switch v := output.(type) {
	case AgentFinish:
		return v, nil
	case AgentAction:
		return a.runTool(nameToToolMap, colorMapping, inputs, []AgentAction{v})
	default:
		return a.runTool(nameToToolMap, colorMapping, inputs, v.([]AgentAction))
	}
}

func (a *AgentExecutor) runTool(nameToToolMap map[string]toolSchema.BaseTool, colorMapping map[string]string, inputs map[string]interface{}, actions []AgentAction) ([]IntermediateStep, error) {
	var result []IntermediateStep
	for _, agentAction := range actions {
		a.CallbackManager.OnAgentAction(agentAction, a.Verbose, "green")
		tool, ok := nameToToolMap[agentAction.Tool]
		if !ok {
			return nil, fmt.Errorf("invalid tool %s", agentAction.Tool)
		}
		color := colorMapping[agentAction.Tool]
		toolRunKwargs := a.agent.(BaseAgent).ToolRunLoggingKwargs()
		observation, err := tool.Run(agentAction.ToolInput, &a.Verbose, color, toolRunKwargs)
		if err != nil {
			return []IntermediateStep{}, err
		}
		result = append(result, IntermediateStep{agentAction, observation})
	}
	return result, nil
}

func (a *AgentExecutor) Call(inputs map[string]interface{}) (map[string]interface{}, error) {
	nameToToolMap := make(map[string]toolSchema.BaseTool)
	var colors []string
	for _, tool := range a.tools {
		nameToToolMap[tool.Name] = tool
		colors = append(colors, tool.Name)
	}
	colorMapping := tools.GetColorMapping(colors, "green")
	var intermediateSteps []IntermediateStep
	iterations := 0
	startTime := time.Now()
	for a.ShouldContinue(iterations, float64(time.Since(startTime).Milliseconds())/1000) {
		nextStepOutput, err := a.TakeNextStep(nameToToolMap, colorMapping, inputs, intermediateSteps)
		if err != nil {
			return nil, err
		}
		switch v := nextStepOutput.(type) {
		case AgentFinish:
			return a.Return(v, intermediateSteps), nil
		default:
			intermediateSteps = append(intermediateSteps, v.([]IntermediateStep)...)
		}
		iterations++
	}
	output := a.agent.(BaseAgent).ReturnStoppedResponse(a.earlyStoppingMethod, intermediateSteps, inputs)
	return a.Return(output, intermediateSteps), nil
}

func (a *AgentExecutor) Return(output AgentFinish, intermediateSteps []IntermediateStep) map[string]interface{} {
	a.CallbackManager.OnAgentFinish(output, "green", a.Verbose)
	finalOutput := output.ReturnValues
	if a.returnIntermediateSteps {
		finalOutput["intermediate_steps"] = intermediateSteps
	}
	return finalOutput
}

func (a *AgentExecutor) GetToolReturn(nextStepOutput IntermediateStep) *AgentFinish {
	nameToToolMap := make(map[string]toolSchema.BaseTool)
	for _, tool := range a.tools {
		nameToToolMap[tool.Name] = tool
	}
	if tool, ok := nameToToolMap[nextStepOutput.Tool]; ok {
		if tool.ReturnDirect {
			//  TODO: not sure what to a
			return &AgentFinish{map[string]interface{}{a.agent.(BaseAgent).ReturnValues()[0]: nextStepOutput.Output}, ""}
		}
	}
	return nil
}
