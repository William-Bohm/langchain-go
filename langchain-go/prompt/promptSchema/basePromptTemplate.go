package promptSchema

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/outputParser/outputParserSchema"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Callable func() string

type BasePromptTemplate struct {
	InputVariables   []string
	OutputParser     outputParserSchema.BaseOutputParser
	PartialVariables map[string]interface{}
	PromptType       string
}

type BasePromptTemplateInterface interface {
	FormatPrompt(kwargs map[string]interface{}) (PromptValue, error) //  abstract
	ValidateVariableNames(values map[string]interface{}) (map[string]interface{}, error)
	Partial(kwargs map[string]interface{}) (BasePromptTemplate, error)
	MergePartialAndUserVariables(kwargs map[string]interface{}) (map[string]interface{}, error)
	Format(kwargs map[string]interface{}) (string, error) //  abastract
	PromptType() string                                   //  abstract
	ToDict(kwargs map[string]interface{}) (map[string]interface{}, error)
	Save(filePath string) error
}

func (bpt *BasePromptTemplate) ValidateVariableNames(values map[string]interface{}) (map[string]interface{}, error) {
	// Implement validation logic
	return values, nil
}

func (bpt *BasePromptTemplate) FormatPrompt(kwargs map[string]interface{}) (PromptValue, error) {

}

func (bpt *BasePromptTemplate) partial(kwargs map[string]interface{}) *BasePromptTemplate {
	promptDict := BasePromptTemplate{
		InputVariables:   make([]string, 0, len(bpt.InputVariables)),
		PartialVariables: make(map[string]interface{}),
		PromptType:       bpt.PromptType,
		OutputParser:     bpt.OutputParser,
	}

	for _, v := range bpt.InputVariables {
		if _, ok := kwargs[v]; !ok {
			promptDict.InputVariables = append(promptDict.InputVariables, v)
		}
	}

	for k, v := range bpt.PartialVariables {
		promptDict.PartialVariables[k] = v
	}

	for k, v := range kwargs {
		promptDict.PartialVariables[k] = v
	}

	return &promptDict
}

func (bpt *BasePromptTemplate) MergePartialAndUserVariables(kwargs map[string]interface{}) map[string]interface{} {
	partialKwargs := make(map[string]interface{})
	for k, v := range bpt.PartialVariables {
		switch v := v.(type) {
		case string:
			partialKwargs[k] = v
		case Callable:
			partialKwargs[k] = v()
		default:
			fmt.Printf("Invalid type %v for value %v\n", reflect.TypeOf(v), v)
		}
	}

	for k, v := range kwargs {
		partialKwargs[k] = v
	}

	return partialKwargs
}
func (bpt *BasePromptTemplate) ToDict(kwargs map[string]interface{}) (map[string]interface{}, error) {
	// Assuming you have a separate method for getting the base dictionary
	// TODO: BELOW!!!!!!!!!!!
	// ******************************* CHECK THE BASEDICT VALUES WITH THE VALUES IN THE MAP TO SEE IF ANY ARE MISSED!!!!!!!!!!!
	baseDict := map[string]interface{}{
		"input_variables":   bpt.InputVariables,
		"output_parser":     bpt.OutputParser,
		"partial_variables": bpt.PartialVariables,
		"prompt_type":       bpt.PromptType,
	}
	if baseDict == nil {
		return nil, errors.New("Error creating base dictionary")
	}
	baseDict["_type"] = bpt.PromptType
	return baseDict, nil
}

func (bpt *BasePromptTemplate) Save(filePath string) error {
	if len(bpt.PartialVariables) != 0 {
		return errors.New("Cannot save prompt with partial variables")
	}

	savePath := filepath.FromSlash(filePath)
	directoryPath := filepath.Dir(savePath)
	err := os.MkdirAll(directoryPath, os.ModePerm)
	if err != nil {
		return err
	}

	promptDict, err := bpt.ToDict(map[string]interface{}{})
	if err != nil {
		return err
	}

	var data []byte
	if strings.HasSuffix(savePath, ".json") {
		data, err = json.MarshalIndent(promptDict, "", "  ")
	} else if strings.HasSuffix(savePath, ".yaml") {
		// Use your favorite YAML library to marshal the data
		// data, err = yaml.Marshal(promptDict)
	} else {
		return errors.New("File must be json or yaml")
	}

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(savePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func NewBasePromptTemplate(inputVars []string, outputParser outputParserSchema.BaseOutputParser, partialVars map[string]interface{}, promptType string) BasePromptTemplate {
	return BasePromptTemplate{
		InputVariables:   inputVars,
		OutputParser:     outputParser,
		PartialVariables: partialVars,
		PromptType:       promptType,
	}
}
