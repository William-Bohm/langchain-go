package llmSchema

import (
	"encoding/json"
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/openai"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func load_llm_from_config(config map[string]interface{}) (BaseLanguageModel, error) {
	// Load BaseLanguageModel from config
	llmType := config["LLMType"].(string)
	switch llmType {
	case "openai":
		return openai.NewFromMap(config)
	default:
		return nil, errors.New("invalid BaseLanguageModel type")
	}
}

func LoadLLM(file string) (BaseLanguageModel, error) {
	// Convert file to absolute path.
	absPath, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	// Determine the file type.
	fileExt := strings.ToLower(filepath.Ext(absPath))
	var config []byte
	switch fileExt {
	case ".json":
		// Load the JSON file.
		config, err = os.ReadFile(absPath)
		if err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		// Load the YAML file.
		config, err = os.ReadFile(absPath)
		if err != nil {
			return nil, err
		}
		var node yaml.Node
		err = yaml.Unmarshal(config, &node)
		if err != nil {
			return nil, err
		}
		config, err = json.Marshal(node)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("File type must be json or yaml")
	}

	// Unmarshal the JSON data into a mapTools[string]interface{}.
	var data map[string]interface{}
	err = json.Unmarshal(config, &data)
	if err != nil {
		return nil, err
	}

	return load_llm_from_config(data)
}
