package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/agent/agentSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/chains"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/tools/mapTools"
	"github.com/William-Bohm/langchain-go/langchain-go/tools/toolSchema"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	defaultRef = os.Getenv("LANGCHAIN_HUB_DEFAULT_REF")
	urlBase    = os.Getenv("LANGCHAIN_HUB_URL_BASE")
	hubPathRe  = regexp.MustCompile(`lc(?P<ref>@[^:]+)?://(?P<path>.*)`)
)

var AGENT_TO_CLASS = map[AgentType]agentSchema.BaseSingleActionAgent{
	zeroShotReactDescription:           &ZeroShotAgent{},
	reactDocstore:                      &ReActDocstoreAgent{},
	selfAskWithSearch:                  &SelfAskWithSearchAgent{},
	conversationalReactDescription:     &ConversationalAgent{},
	chatZeroShotReactDescription:       &ChatAgent{},
	chatConversationalReactDescription: &ConversationalChatAgent{},
}

const URL_BASE = "https://raw.githubusercontent.com/hwchase17/langchain-hub/master/agents/"

func LoadAgentFromTools(config map[string]interface{}, llm llmSchema.BaseLLM, tools []toolSchema.BaseTool, kwargs map[string]interface{}) (*agentSchema.BaseSingleActionAgent, error) {
	configType := config["_type"].(AgentType)
	delete(config, "_type")
	if _, ok := AGENT_TO_CLASS[configType]; !ok {
		return &agentSchema.BaseSingleActionAgent{}, errors.New("Loading " + string(configType) + " agent not supported")
	}
	agentClass := AGENT_TO_CLASS[configType]
	combinedConfig := mapTools.MergeMaps(config, kwargs)
	return agentClass.FromLLMAndTools(llm, tools, combinedConfig["callbackManger"].(callbackSchema.BaseCallbackManager), combinedConfig)
}

func LoadAgentFromConfig(config map[string]interface{}, llm llmSchema.BaseLLM, tools []toolSchema.BaseTool, kwargs map[string]interface{}) (*agentSchema.BaseSingleActionAgent, error) {
	var err error
	if _, ok := config["_type"]; !ok {
		return &agentSchema.BaseSingleActionAgent{}, errors.New("Must specify an agent Type in config")
	}
	loadFromTools := config["load_from_llm_and_tools"].(bool)
	delete(config, "load_from_llm_and_tools")
	if loadFromTools {
		if llm.Id == "" && llm.CallbackManager == "" {
			return &agentSchema.BaseSingleActionAgent{}, errors.New("If `load_from_llm_and_tools` is set to True, then LLM must be provided. Make sure the LLM has an assigned ID and/or CallbackManger!")
		}
		if tools == nil {
			return &agentSchema.BaseSingleActionAgent{}, errors.New("If `load_from_llm_and_tools` is set to True, then tools must be provided")
		}
		return LoadAgentFromTools(config, llm, tools, kwargs)
	}
	configType := config["_type"].(AgentType)
	delete(config, "_type")
	if _, ok := AGENT_TO_CLASS[configType]; !ok {
		return &agentSchema.BaseSingleActionAgent{}, errors.New("Loading " + string(configType) + " agent not supported")
	}
	agentClass := AGENT_TO_CLASS[configType]
	if _, ok := config["llm_chain"]; ok {
		config["llm_chain"], err = chains.LoadChainFromConfig(config["llm_chain"].(map[string]interface{}), kwargs)
		if err != nil {
			return &agentSchema.BaseSingleActionAgent{}, err
		}
	} else if _, ok := config["llm_chain_path"]; ok {
		config["llm_chain"], err = chains.LoadChain(config["llm_chain_path"].(string), kwargs)
	} else {
		return &agentSchema.BaseSingleActionAgent{}, errors.New("One of `llm_chain` and `llm_chain_path` should be specified.")
	}
	combinedConfig := mapTools.MergeMaps(config, kwargs)
	return agentClass.FromConfig(combinedConfig), nil
}

func LoadAgent(path string, kwargs map[string]interface{}) (*agentSchema.BaseSingleActionAgent, error) {
	if hubResult, err := LoadAgentFromHub(path, LoadAgentFromFile, "agents", []string{"json", "yaml"}, map[string]interface{}{}); hubResult != nil {
		if err != nil {
			return &agentSchema.BaseSingleActionAgent{}, err
		}
		return hubResult, nil
	} else {
		return LoadAgentFromFile(path, kwargs)
	}
}

func LoadAgentFromFile(file string, kwargs map[string]interface{}) (*agentSchema.BaseSingleActionAgent, error) {
	var config map[string]interface{}
	filePath := path.Ext(file)
	if filePath == ".json" {
		fileBytes, _ := os.ReadFile(file)
		err := json.Unmarshal(fileBytes, &config)
		if err != nil {
			return nil, err
		}
	} else if filePath == ".yaml" {
		fileBytes, _ := os.ReadFile(file)
		err := yaml.Unmarshal(fileBytes, &config)
		if err != nil {
			return nil, err
		}
	} else {
		return &agentSchema.BaseSingleActionAgent{}, errors.New("File type must be json or yaml")
	}
	return LoadAgentFromConfig(config, llmSchema.BaseLLM{}, []toolSchema.BaseTool{}, kwargs)
}

func LoadAgentFromHub(
	path string,
	loader func(string, map[string]interface{}) (*agentSchema.BaseSingleActionAgent, error),
	validPrefix string,
	validSuffixes []string,
	kwargs map[string]interface{},
) (*agentSchema.BaseSingleActionAgent, error) {
	if _, err := url.ParseRequestURI(path); err != nil || !hubPathRe.MatchString(path) {
		return &agentSchema.BaseSingleActionAgent{}, nil
	}

	matches := hubPathRe.FindStringSubmatch(path)
	ref := matches[1]
	if ref == "" {
		ref = defaultRef
	} else {
		ref = ref[1:]
	}
	remotePathStr := matches[2]
	remotePath := filepath.Clean(remotePathStr)
	if strings.Split(remotePath, "/")[0] != validPrefix {
		return &agentSchema.BaseSingleActionAgent{}, nil
	}
	if !contains(validSuffixes, filepath.Ext(remotePath)) {
		return &agentSchema.BaseSingleActionAgent{}, fmt.Errorf("Unsupported file type.")
	}

	fullURL := urlBase + "/" + ref + "/" + remotePath

	resp, err := http.Get(fullURL)
	if err != nil || resp.StatusCode != 200 {
		return &agentSchema.BaseSingleActionAgent{}, fmt.Errorf("Could not find file at %s", fullURL)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)

	tmpDirName := os.TempDir() + "/" + uuid.New().String()
	err = os.Mkdir(tmpDirName, 0755)
	if err != nil {
		return nil, err
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {

		}
	}(tmpDirName)

	file := tmpDirName + "/" + filepath.Base(remotePath)
	err = os.WriteFile(file, body, 0644)
	if err != nil {
		return nil, err
	}

	return loader(file, kwargs)
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
