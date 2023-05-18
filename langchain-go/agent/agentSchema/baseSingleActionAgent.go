package agentSchema

import (
	"encoding/json"
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/tools/toolSchema"
	"os"
	"path/filepath"
)

type SingleActionAgent interface {
	Plan(intermediateSteps []IntermediateStep, kwargs map[string]interface{}) (AgentAction, error)
	InputKeys() []string
	Dict(kwargs map[string]interface{}) map[string]interface{}
}

type BaseSingleActionAgent struct {
	SingleActionAgent
	BaseAgent
}

func (b *BaseSingleActionAgent) FromConfig(config map[string]interface{}) *BaseSingleActionAgent {
	//  abstract class
	return &BaseSingleActionAgent{}
}

func (b *BaseSingleActionAgent) ReturnValues() []string {
	return []string{"output"}
}

func (b *BaseSingleActionAgent) GetAllowedTools() []string {
	return nil
}

func (b *BaseSingleActionAgent) ReturnStoppedResponse(earlyStoppingMethod string, intermediateSteps []IntermediateStep, kwargs map[string]interface{}) (AgentFinish, error) {
	if earlyStoppingMethod == "force" {
		agentReturnValues := make(map[string]interface{})
		agentReturnValues["output"] = "BaseAgent stopped due to iteration limit or time limit."
		agentReturnValues["details"] = ""
		return AgentFinish{agentReturnValues, "BaseAgent stopped due to iteration limit or time limit."}, nil
	} else {
		return AgentFinish{}, errors.New("Got unsupported early_stopping_method `" + earlyStoppingMethod + "`")
	}
}

func (b *BaseSingleActionAgent) FromLLMAndTools(llm llmSchema.BaseLanguageModel, tools []toolSchema.BaseTool, callbackManger callbackSchema.BaseCallbackManager, kwargs map[string]interface{}) (*BaseSingleActionAgent, error) {
	return nil, errors.New("not implemented")
}

func (b *BaseSingleActionAgent) AgentType() string {
	return "BaseSingleActionAgent"
}

func (b *BaseSingleActionAgent) Save(filePath string) error {
	savePath := filePath
	directoryPath := filepath.Dir(savePath)
	os.MkdirAll(directoryPath, os.ModePerm)

	agentDict := b.Dict(map[string]interface{}{})

	var err error
	fileExt := filepath.Ext(filePath)
	if fileExt == ".json" {
		jsonBytes, err := json.MarshalIndent(agentDict, "", "    ")
		if err != nil {
			return err
		}
		err = os.WriteFile(filePath, jsonBytes, 0644)
	} else if fileExt == ".yaml" {
		return errors.New(filePath + " must be json or yaml")
	} else {
		return errors.New(filePath + " must be json or yaml")
	}

	return err
}

func (b *BaseSingleActionAgent) ToolRunLoggingKwargs() map[string]interface{} {
	return map[string]interface{}{}
}

func NewBaseSingleActionAgent() *BaseSingleActionAgent {
	return &BaseSingleActionAgent{}
}
