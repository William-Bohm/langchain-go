package agentSchema

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type BaseMultiActionAgent struct {
	BaseAgent
}

type BaseAgent interface {
	ReturnValues() []string
	GetAllowedTools() []string
	Plan(intermediateSteps []IntermediateStep, kwargs map[string]interface{}) (AgentAction, error)
	InputKeys() []string
	ReturnStoppedResponse(earlyStoppingMethod string, intermediateSteps []IntermediateStep, kwargs map[string]interface{}) (AgentFinish, error)
	AgentType() string
	Dict(kwargs map[string]interface{}) map[string]interface{}
	Save(filePath string) error
	ToolRunLoggingKwargs() map[string]interface{}
}

func (agent *BaseMultiActionAgent) ReturnValues() []string {
	return []string{"output"}
}

func (agent *BaseMultiActionAgent) GetAllowedTools() []string {
	return nil
}

func (agent *BaseMultiActionAgent) ReturnStoppedResponse(earlyStoppingMethod string, intermediateSteps []struct {
	AgentAction
	string
}, kwargs map[string]interface{}) AgentFinish {

	if earlyStoppingMethod == "force" {
		return AgentFinish{map[string]interface{}{"intermediateSteps": intermediateSteps}, "BaseAgent stopped due to max iterations."}
	} else {
		panic(errors.New("Got unsupported early_stopping_method `" + earlyStoppingMethod + "`"))
	}
}

func (agent *BaseMultiActionAgent) Dict(kwargs map[string]interface{}) map[string]interface{} {
	dict := make(map[string]interface{})
	dict["_type"] = agent.AgentType()
	return dict
}

func (agent *BaseMultiActionAgent) Save(filePath string) error {
	savePath := filePath
	directoryPath := filepath.Dir(savePath)

	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		os.MkdirAll(directoryPath, os.ModePerm)
	}

	agentDict := agent.Dict(nil)

	if strings.HasSuffix(savePath, ".json") {
		file, _ := json.MarshalIndent(agentDict, "", " ")
		_ = os.WriteFile(filePath, file, 0644)
	} else if strings.HasSuffix(savePath, ".yaml") {
		// TODO: yaml handling code goes here
	} else {
		panic(errors.New(savePath + " must be json or yaml"))
	}
	return nil
}

func (agent *BaseMultiActionAgent) ToolRunLoggingKwargs() map[string]interface{} {
	return make(map[string]interface{})
}

// Here are the remaining methods that were not defined

func (agent *BaseMultiActionAgent) Plan(intermediateSteps []struct {
	AgentAction
	string
}, kwargs map[string]interface{}) ([]AgentAction, AgentFinish) {
	// Implement this method in the respective agents
	return nil, AgentFinish{}
}

func (agent *BaseMultiActionAgent) InputKeys() []string {
	// Implement this method in the respective agents
	return nil
}

func (agent *BaseMultiActionAgent) AgentType() string {
	// Implement this method in the respective agents
	return ""
}
