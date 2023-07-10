package tracer

/*
The session logic towards the end of some functions in this file may be completley wrong, did not think very much about that code
*/

import (
	"bytes"
	"encoding/json"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/tracer/tracerSchema"
	"log"
	"net/http"
	"os"
)

type BaseLangChainTracer struct {
	*tracerSchema.BaseTracer
	alwaysVerbose bool
	endpoint      string
	headers       map[string]string
}

func NewBaseLangChainTracer() *BaseLangChainTracer {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	apiKey, exists := os.LookupEnv("LANGCHAIN_API_KEY")
	if exists {
		headers["x-api-key"] = apiKey
	}

	endpoint, exists := os.LookupEnv("LANGCHAIN_ENDPOINT")
	if !exists {
		endpoint = "http://localhost:8000"
	}

	baseTracer := tracerSchema.NewBaseTracer("baseLangchainTracer", 1, nil)

	return &BaseLangChainTracer{
		BaseTracer:    baseTracer,
		alwaysVerbose: true,
		endpoint:      endpoint,
		headers:       headers,
	}
}

func (bt *BaseLangChainTracer) PersistRun(run interface{}) {
	var endpoint string
	// Here we need to use Go equivalent of isinstance() check and decide endpoint
	// based on the actual type of `run`.
	// Since the actual type is not provided, using a placeholder switch statement
	switch run.(type) {
	case string: // Placeholder for LLMRun
		endpoint = bt.endpoint + "/llm-runs"
	case int: // Placeholder for ChainRun
		endpoint = bt.endpoint + "/chain-runs"
	default: // Placeholder for ToolRun
		endpoint = bt.endpoint + "/tool-runs"
	}

	data, err := json.Marshal(run)
	if err != nil {
		log.Fatalf("Failed to marshal run: %v", err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
		return
	}

	for key, value := range bt.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to persist run: %v", err)
		return
	}
	defer resp.Body.Close()
}

func (bt *BaseLangChainTracer) PersistSession(sessionCreate interface{}) (interface{}, error) {
	endpoint := bt.endpoint + "/sessions"

	data, err := json.Marshal(sessionCreate)
	if err != nil {
		log.Printf("Failed to marshal sessionCreate: %v", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	for key, value := range bt.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to create session, using default session: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(bt.Session)
	if err != nil {
		log.Printf("Failed to decode response body: %v", err)
		return nil, err
	}

	return bt.Session, nil
}

func (bt *BaseLangChainTracer) LoadSession(sessionName string) (interface{}, error) {
	endpoint := bt.endpoint + "/sessions?name=" + sessionName

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	for key, value := range bt.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to load session %s, using empty session: %v", sessionName, err)
		return nil, err
	}
	defer resp.Body.Close()

	tracerSession := tracerSchema.NewTracerSession(sessionName, map[string]interface{}{})
	err = json.NewDecoder(resp.Body).Decode(&tracerSession)
	if err != nil {
		log.Printf("Failed to decode response body: %v", err)
		return nil, err
	}

	bt.Session = tracerSession

	return tracerSession, nil
}
