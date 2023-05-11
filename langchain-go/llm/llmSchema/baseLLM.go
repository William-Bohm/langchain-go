package llmSchema

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/config/defaults"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// BaseLanguageModel Abstract methods all BaseLanguageModel openaiClient's should define
type BaseLanguageModel interface {
	Generate(prompts []string, stop string) (LLMResult, error)
	GetNumTokensFromMessage(messages []rootSchema.BaseMessageInterface) (int, error)
	GetNumTokensFromText(text string) (int, error)
}

type LLMResult struct {
	Generations [][]Generation
	LLMOutput   map[string]interface{}
}

type Generation struct {
	Text           string                 `json:"text"`                      // Generated text output
	GenerationInfo map[string]interface{} `json:"generation_info,omitempty"` // Raw generation info response from the provider
}

// Reduct boilerplate for each BaseLanguageModel by defining functions they will all use
type BaseLLM struct {
	BaseLanguageModel
	Id              string
	Verbose         bool
	CountTokens     bool
	CallbackManager string
	Cache           string
	LLMType         string
}

func NewDefaultBaseLLM(LLMType string, apiKey string) *BaseLLM {
	return &BaseLLM{
		Id:              uuid.New().String(),
		Verbose:         defaults.GetDefaultVerbose(),
		CountTokens:     defaults.GetDefaultCountTokens(),
		CallbackManager: defaults.GetDefaultCallbackManager(),
		Cache:           defaults.GetDefaultCache(),
		LLMType:         LLMType,
	}
}

func NewBaseLLM(attrs map[string]interface{}, LLMType string) (*BaseLLM, error) {
	// create a new BaseLLM from a mapTools of attributes
	baseLLM := &BaseLLM{}
	if val, ok := attrs["Verbose"]; ok {
		baseLLM.Verbose = val.(bool)
	} else {
		baseLLM.Verbose = defaults.GetDefaultVerbose()
	}

	if val, ok := attrs["CallbackManager"]; ok {
		baseLLM.CallbackManager = val.(string)
	} else {
		baseLLM.CallbackManager = defaults.GetDefaultCallbackManager()
	}

	if val, ok := attrs["Cache"]; ok {
		baseLLM.Cache = val.(string)
	} else {
		baseLLM.Cache = defaults.GetDefaultCache()
	}

	if val, ok := attrs["LLMType"]; ok {
		baseLLM.LLMType = val.(string)
	} else {
		baseLLM.LLMType = LLMType
	}

	if val, ok := attrs["Id"]; ok {
		baseLLM.Id = val.(string)
	} else {
		baseLLM.Id = uuid.New().String()
	}

	return baseLLM, nil
}

// TODO: Implement these functions

func (llm *BaseLLM) SetVerbose(verbose bool) {
	llm.Verbose = verbose
}

func (llm *BaseLLM) SetCallBackHandler(handler string) {
	llm.CallbackManager = handler
}

func (llm *BaseLLM) SetCache(cache string) {
	llm.Cache = cache
}

func (llm *BaseLLM) SetLLMType(llmType string) {
	llm.LLMType = llmType
}

func (llm *BaseLLM) GetLLMType() string {
	return llm.LLMType
}

func (llm *BaseLLM) SaveLLMToFile(filePath string) error {
	savePath := filepath.Clean(filePath)

	// Create the directory if it does not exist
	dirPath := filepath.Dir(savePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	// Get the mapTools representation of the object
	llmMap := llm.baseLLMToMap()

	// Determine the file format based on the file extension
	fileExt := filepath.Ext(savePath)

	var err error
	var data []byte
	switch strings.ToLower(fileExt) {
	case ".json":
		data, err = json.MarshalIndent(llmMap, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(llmMap)
	default:
		return errors.New("file must be JSON or YAML")
	}

	if err != nil {
		return err
	}

	// Write the data to the file
	err = os.WriteFile(savePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error marshaling YAML/JSON data: %v", err)
	}

	return nil
}
func (llm *BaseLLM) baseLLMToMap() map[string]interface{} {
	// Convert BaseLLM struct fields to a mapTools
	return map[string]interface{}{
		"Id":              llm.Id,
		"Verbose":         llm.Verbose,
		"CountTokens":     llm.CountTokens,
		"CallbackManager": llm.CallbackManager,
		"Cache":           llm.Cache,
		"LLMType":         llm.LLMType,
	}
}
