package memory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/William-Bohm/langchain-go/langchain-go/memory/memorySchema"
)

type MotorheadMemory struct {
	memorySchema.BaseChatMemory
	URL       string
	Timeout   time.Duration
	MemoryKey string
	SessionID string
	Context   *string
}

func NewMotorheadMemory(sessionID string, baseChatMemory memorySchema.BaseChatMemory) *MotorheadMemory {
	return &MotorheadMemory{
		BaseChatMemory: baseChatMemory,
		URL:            "http://localhost:8080",
		Timeout:        3000 * time.Millisecond,
		MemoryKey:      "history",
		SessionID:      sessionID,
	}
}

func (m *MotorheadMemory) Init() error {
	client := &http.Client{Timeout: m.Timeout}
	url := fmt.Sprintf("%s/sessions/%s/memory", m.URL, m.SessionID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resData map[string]interface{}
	err = json.Unmarshal(body, &resData)
	if err != nil {
		return err
	}

	messages := resData["messages"].([]interface{})
	context := resData["context"].(string)

	for _, msg := range messages {
		message := msg.(map[string]interface{})
		role := message["role"].(string)
		content := message["content"].(string)

		if role == "AI" {
			m.ChatMemory.AddAIMessage(content)
		} else {
			m.ChatMemory.AddUserMessage(content)
		}
	}

	if context != "" && context != "NONE" {
		m.Context = &context
	}

	return nil
}

func (m *MotorheadMemory) LoadMemoryVariables(values map[string]interface{}) (map[string]interface{}, error) {
	if m.BaseChatMemory.ReturnMessages {
		return map[string]interface{}{m.MemoryKey: m.ChatMemory.Messages()}, nil
	} else {
		messages, err := m.ChatMemory.Messages()
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{m.MemoryKey: rootSchema.GetBufferString(messages)}, nil
	}
}

func (m *MotorheadMemory) MemoryVariables() []string {
	return []string{m.MemoryKey}
}

func (m *MotorheadMemory) SaveContext(inputs map[string]interface{}, outputs map[string]string) error {
	inputStr, outputStr, err := m.GetInputOutput(inputs, outputs)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: m.Timeout}
	url := fmt.Sprintf("%s/sessions/%s/memory", m.URL, m.SessionID)
	postData := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "Human", "content": inputStr},
			{"role": "AI", "content": outputStr},
		},
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return m.BaseChatMemory.SaveContext(inputs, outputs)
}

func (m *MotorheadMemory) _get_input_output(inputs map[string]interface{}, outputs map[string]string) (string, string) {
	var inputStr, outputStr string
	for k, v := range inputs {
		if _, ok := outputs[k]; ok {
			inputStr = fmt.Sprintf("%v", v)
			outputStr = outputs[k]
			break
		}
	}

	return inputStr, outputStr
}
