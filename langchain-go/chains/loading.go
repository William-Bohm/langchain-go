package chains

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ChainLoader func(config map[string]interface{}, kwargs map[string]interface{}) Chain

var typeToLoaderDict = map[string]ChainLoader{
	"api_chain":                       loadAPIChain,
	"hyde_chain":                      loadHydeChain,
	"llm_chain":                       loadLLMChain,
	"llm_bash_chain":                  loadLLMBashChain,
	"llm_checker_chain":               loadLLMCheckerChain,
	"llm_math_chain":                  loadLLMMathChain,
	"llm_requests_chain":              loadLLMRequestsChain,
	"pal_chain":                       loadPalChain,
	"qa_with_sources_chain":           loadQAWithSourcesChain,
	"stuff_documents_chain":           loadStuffDocumentsChain,
	"map_reduce_documents_chain":      loadMapReduceDocumentsChain,
	"map_rerank_documents_chain":      loadMapRerankDocumentsChain,
	"refine_documents_chain":          loadRefineDocumentsChain,
	"sql_database_chain":              loadSQLDatabaseChain,
	"vector_db_qa_with_sources_chain": loadVectorDBQAWithSourcesChain,
	"vector_db_qa":                    loadVectorDBQA,
}

var (
	defaultRef = os.Getenv("LANGCHAIN_HUB_DEFAULT_REF")
	urlBase    = os.Getenv("LANGCHAIN_HUB_URL_BASE")
	hubPathRe  = regexp.MustCompile(`lc(?P<ref>@[^:]+)?://(?P<path>.*)`)
)

func LoadChainFromConfig(config map[string]interface{}, kwargs map[string]interface{}) (Chain, error) {
	if _, ok := config["_type"]; !ok {
		return nil, errors.New("Must specify a chain Type in config")
	}
	configType := config["_type"].(string)

	if _, ok := typeToLoaderDict[configType]; !ok {
		return nil, errors.New("Loading " + configType + " chain not supported")
	}

	chainLoader := typeToLoaderDict[configType]
	return chainLoader(config, kwargs), nil
}

func LoadChain(path string, kwargs map[string]interface{}) (Chain, error) {
	if hubResult, err := LoadChainFromHub(path, LoadChainFromFile, "chains", []string{"json", "yaml"}, kwargs); err == nil {
		return hubResult, nil
	} else {
		return LoadChainFromFile(path, kwargs)
	}
}

func LoadChainFromFile(file string, kwargs map[string]interface{}) (Chain, error) {
	var config map[string]interface{}

	fileData, _ := ioutil.ReadFile(file)
	fileType := filepath.Ext(file)
	if fileType == ".json" {
		json.Unmarshal(fileData, &config)
	} else if fileType == ".yaml" {
		yaml.Unmarshal(fileData, &config)
	} else {
		return nil, errors.New("File type must be json or yaml")
	}

	if verbose, ok := kwargs["verbose"]; ok {
		config["verbose"] = verbose
		delete(kwargs, "verbose")
	}
	if memory, ok := kwargs["memory"]; ok {
		config["memory"] = memory
		delete(kwargs, "memory")
	}

	return LoadChainFromConfig(config, kwargs)
}

func LoadChainFromHub(
	path string,
	loader func(string, map[string]interface{}) (Chain, error),
	validPrefix string,
	validSuffixes []string,
	kwargs map[string]interface{},
) (Chain, error) {
	if _, err := url.ParseRequestURI(path); err != nil || !hubPathRe.MatchString(path) {
		return nil, nil
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
		return nil, nil
	}
	if !contains(validSuffixes, filepath.Ext(remotePath)) {
		return nil, fmt.Errorf("Unsupported file type.")
	}

	fullURL := urlBase + "/" + ref + "/" + remotePath

	resp, err := http.Get(fullURL)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("Could not find file at %s", fullURL)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	tmpDirName := os.TempDir() + "/" + uuid.New().String()
	os.Mkdir(tmpDirName, 0755)
	defer os.RemoveAll(tmpDirName)

	file := tmpDirName + "/" + filepath.Base(remotePath)
	os.WriteFile(file, body, 0644)

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
