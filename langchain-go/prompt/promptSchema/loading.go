package promptSchema

import (
	"encoding/json"
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/outputParser/outputParserSchema"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func loadPromptFromFile(file string) (BasePromptTemplateInterface, error) {
	filePath := filepath.Clean(file)
	fileExt := filepath.Ext(filePath)

	var config map[string]interface{}

	switch fileExt {
	case ".json":
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(fileContent, &config)
		if err != nil {
			return nil, err
		}
	case ".yaml":
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(fileContent, &config)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported file type: " + fileExt)
	}

	return loadPromptFromConfig(config)
}

func loadPromptFromConfig(config map[string]interface{}) (BasePromptTemplateInterface, error) {
	// once the config has been aquired load the template for it
	configType, ok := config["_type"]
	if !ok {
		configType = "prompt"
	}

	delete(config, "_type")

	if configType.(string) == "prompt" {
		template, err := loadPromptTemplate(config)
		if err != nil {
			return nil, err
		}
		return &template, nil
	} else if configType.(string) == "few_shot" {
		template, err := loadFewShotPrompt(config)
		if err != nil {
			return nil, err
		}
		return &template, nil
	} else {
		return nil, errors.New("loading " + configType.(string) + " prompt not supported")

	}
}

func loadFewShotPrompt(config map[string]interface{}) (FewShotPromptTemplate, error) {
	var err error
	// Load the suffix and prefix templates.
	config, err = loadTemplate("suffix", config)
	if err != nil {
		return FewShotPromptTemplate{}, err
	}
	config, err = loadTemplate("prefix", config)
	if err != nil {
		return FewShotPromptTemplate{}, err
	}

	// Load the example prompt.
	if examplePromptPath, ok := config["example_prompt_path"]; ok {
		if _, ok := config["example_prompt"]; ok {
			return FewShotPromptTemplate{}, errors.New("only one of example_prompt and example_prompt_path should be specified")
		}
		config["example_prompt"], err = loadPrompt(examplePromptPath.(string))
		if err != nil {
			return FewShotPromptTemplate{}, err
		}
		delete(config, "example_prompt_path")
	} else {
		config["example_prompt"], err = loadPromptFromConfig(config["example_prompt"].(map[string]interface{}))
		if err != nil {
			return FewShotPromptTemplate{}, err
		}
	}

	// Load the examples.
	config, err = loadExamples(config)
	if err != nil {
		return FewShotPromptTemplate{}, err
	}
	config, err = loadOutputParser(config)
	if err != nil {
		return FewShotPromptTemplate{}, err
	}
	// TODO: make the config settings actually set setting for the template
	return NewFewShotPromptTemplate(config), nil
}

func loadPromptTemplate(config map[string]interface{}) (BasePromptTemplate, error) {
	var err error
	// Load the template from disk if necessary.
	config, err = loadTemplate("template", config)
	if err != nil {
		return BasePromptTemplate{}, err
	}
	config, err = loadOutputParser(config)
	if err != nil {
		return BasePromptTemplate{}, err
	}
	return newPromptTemplate(config), nil
}

func newPromptTemplate(config map[string]interface{}) BasePromptTemplate {
	template := BasePromptTemplate{
		InputVariables:   config["inputVariables"].([]string),
		OutputParser:     config["outputParser"].(outputParserSchema.BaseOutputParser),
		PartialVariables: config["partialVariables"].(map[string]interface{}),
		PromptType:       config["promptType"].(string),
	}

	return template
}

func loadPrompt(path string) (BasePromptTemplateInterface, error) {
	// Check if the path is from LangChainHub or local file system.
	// tryLoadFromHub() function should be implemented.
	if hubResult, err := tryLoadFromHub(path, loadPromptFromFile, "prompt", []string{"py", "json", "yaml"}); err == nil {
		return hubResult, nil
	} else {
		return loadPromptFromFile(path)
	}
}

func loadTemplate(varName string, config map[string]interface{}) (map[string]interface{}, error) {
	varNamePath := varName + "_path"

	if templatePath, ok := config[varNamePath]; ok {
		if _, ok := config[varName]; ok {
			return nil, errors.New("both `" + varNamePath + "` and `" + varName + "` cannot be provided")
		}

		path := templatePath.(string)
		ext := filepath.Ext(path)

		var content []byte
		var err error
		if ext == ".txt" {
			content, err = ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("invalid file format")
		}

		config[varName] = string(content)
		delete(config, varNamePath)
	}

	return config, nil
}

func loadExamples(config map[string]interface{}) (map[string]interface{}, error) {
	examples, ok := config["examples"]
	if !ok {
		return nil, errors.New("examples not found in config")
	}

	switch examples.(type) {
	case []interface{}:
		break
	case string:
		filePath := examples.(string)
		var content []byte
		var err error

		content, err = ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		ext := filepath.Ext(filePath)

		if ext == ".json" {
			var jsonData []map[string]interface{}
			err = json.Unmarshal(content, &jsonData)
			if err != nil {
				return nil, err
			}
			config["examples"] = jsonData
		} else if strings.ToLower(ext) == ".yaml" || strings.ToLower(ext) == ".yml" {
			// Implement YAML loading here, using your preferred library.
		} else {
			return nil, errors.New("invalid file format. Only json or yaml formats are supported")
		}
	default:
		return nil, errors.New("invalid examples format. Only list or string are supported")
	}

	return config, nil
}

func loadOutputParser(config map[string]interface{}) (map[string]interface{}, error) {
	if outputParsers, ok := config["output_parsers"]; ok {
		if outputParsers != nil {
			configData := outputParsers.(map[string]interface{})
			outputParserType := configData["_type"].(string)
			if outputParserType == "regex_parser" {
				outputParser := NewRegexParser(configData)
				config["output_parsers"] = outputParser
			} else {
				return nil, errors.New("unsupported output parser " + outputParserType)
			}
		}
	}

	return config, nil
}
