package agentSchema

import (
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/chains"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/memory/memorySchema"
	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/tools/toolSchema"
	"strings"
)

type Agent struct {
	BaseSingleActionAgent
	llmChain     chains.LLMChain
	outputParser AgentOutputParser
	allowedTools []string
}

func (a *Agent) GetAllowedTools() []string {
	return a.allowedTools
}

func (a *Agent) ReturnValues() []string {
	return []string{"output"}
}

func (a *Agent) _FixText(text string) (string, error) {
	return "", errors.New("fix_text not implemented for this agent.")
}

func (a *Agent) _Stop() []string {
	return []string{fmt.Sprintf("\n%s", strings.TrimSpace(a.ObservationPrefix())),
		fmt.Sprintf("\n\t%s", strings.TrimSpace(a.ObservationPrefix()))}
}

func (a *Agent) _ConstructScratchpad(intermediateSteps []AgentAction) string {
	thoughts := ""
	for _, step := range intermediateSteps {
		thoughts += step.Log
		thoughts += fmt.Sprintf("\n%s\n%s", a.ObservationPrefix(), a.LLMPrefix())
	}
	return thoughts
}

func (a *Agent) Plan(intermediateSteps []AgentAction, kwargs map[string]interface{}) (AgentAction, error) {
	fullInputs := a.GetFullInputs(intermediateSteps, kwargs)
	fullOutput, err := a.llmChain.Predict(fullInputs)
	if err != nil {
		return AgentAction{}, err
	}
	return a.outputParser.Parse(fullOutput).(AgentAction), nil
}

func (a *Agent) GetFullInputs(intermediateSteps []AgentAction, kwargs map[string]interface{}) map[string]interface{} {
	thoughts := a._ConstructScratchpad(intermediateSteps)
	newInputs := map[string]interface{}{
		"agent_scratchpad": thoughts,
		"stop":             a._Stop(),
	}
	fullInputs := kwargs
	for k, v := range newInputs {
		fullInputs[k] = v
	}
	return fullInputs
}

func (a *Agent) InputKeys() []string {
	// TODO: Convert `set` operation into Go equivalent
	return a.llmChain.InputKeys()
}

func (a *Agent) ObservationPrefix() string {
	// To be implemented by the concrete implementation
	panic("not implemented")
}

func (a *Agent) LLMPrefix() string {
	// To be implemented by the concrete implementation
	panic("not implemented")
}

func (a *Agent) ValidatePrompt(values map[string]interface{}) map[string]interface{} {
	// TODO: Complete this method according to the Python code
	return values
}

func (a *Agent) CreatePrompt(tools []toolSchema.BaseTool) promptSchema.BasePromptTemplate {
	// To be implemented by the concrete implementation
	panic("not implemented")
}

func (a *Agent) _ValidateTools(tools []toolSchema.BaseTool) {
	// To be implemented by the concrete implementation
	panic("not implemented")
}

func (a *Agent) GetDefaultOutputParser(kwargs map[string]interface{}) *AgentOutputParser {
	// To be implemented by the concrete implementation
	panic("not implemented")
}

func (a *Agent) FromLLMAndTools(llm llmSchema.BaseLanguageModel, tools []toolSchema.BaseTool, callbackManager callbackSchema.BaseCallbackManager, outputParser *AgentOutputParser, kwargs map[string]interface{}) Agent {
	a._ValidateTools(tools)

	llmChain := chains.LLMChain{
		BaseChain: chains.BaseChain{
			Memory:          kwargs["memory"].(memorySchema.BaseMemory),
			CallbackManager: callbackManager,
			Verbose:         false,
		},
		LLM:       llm,
		Prompt:    a.CreatePrompt(tools),
		OutputKey: "text",
	}
	toolNames := make([]string, len(tools))
	for i, tool := range tools {
		toolNames[i] = tool.Name
	}
	_outputParser := outputParser
	if _outputParser == nil {
		_outputParser = a.GetDefaultOutputParser(kwargs)
	}
	return Agent{
		llmChain:     llmChain,
		allowedTools: toolNames,
		outputParser: *_outputParser,
	}

}

func (a *Agent) ReturnStoppedResponse(earlyStoppingMethod string, intermediateSteps []AgentAction, kwargs map[string]interface{}) (AgentFinish, error) {
	if earlyStoppingMethod == "force" {
		return AgentFinish{
			ReturnValues: map[string]interface{}{"output": "BaseAgent stopped due to iteration limit or time limit."},
			Log:          "",
		}, nil
	} else if earlyStoppingMethod == "generate" {
		thoughts := ""
		for _, step := range intermediateSteps {
			thoughts += step.Log
			thoughts += fmt.Sprintf("\n%s\n%s", a.ObservationPrefix(), a.LLMPrefix())
		}
		thoughts += "\n\nI now need to return a final answer based on the previous steps:"
		newInputs := map[string]interface{}{
			"agent_scratchpad": thoughts,
			"stop":             a._Stop(),
		}
		fullInputs := kwargs
		for k, v := range newInputs {
			fullInputs[k] = v
		}
		fullOutput, err := a.llmChain.Predict(fullInputs)
		if err != nil {
			return AgentFinish{}, err
		}
		parsedOutput := a.outputParser.Parse(fullOutput)
		if _, ok := parsedOutput.(AgentFinish); ok {
			return parsedOutput.(AgentFinish), nil
		} else {
			return AgentFinish{
				ReturnValues: map[string]interface{}{"output": fullOutput},
				Log:          fullOutput,
			}, nil
		}
	} else {
		return AgentFinish{}, errors.New("early_stopping_method should be one of force or generate")
	}
}

func (a *Agent) ToolRunLoggingKwargs() map[string]string {
	return map[string]string{
		"llm_prefix":         a.LLMPrefix(),
		"observation_prefix": a.ObservationPrefix(),
	}
}
