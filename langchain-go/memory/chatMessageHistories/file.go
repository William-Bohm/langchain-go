package chatMessageHistories

import (
	"encoding/json"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
	"io/ioutil"
	"os"
)

type FileChatMessageHistory struct {
	filePath string
}

func NewFileChatMessageHistory(filePath string) (*FileChatMessageHistory, error) {
	f := &FileChatMessageHistory{filePath: filePath}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := ioutil.WriteFile(filePath, []byte("[]"), 0644)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *FileChatMessageHistory) Messages() ([]rootSchema.BaseMessageInterface, error) {
	data, err := ioutil.ReadFile(f.filePath)
	if err != nil {
		return nil, err
	}

	var messages []rootSchema.BaseMessageInterface
	err = json.Unmarshal(data, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (f *FileChatMessageHistory) AddUserMessage(message string) error {
	humanMessage := rootSchema.NewHumanMessage(message)
	return f.append(humanMessage)
}

func (f *FileChatMessageHistory) AddAIMessage(message string) error {
	aiMessage := rootSchema.NewAIMessage(message)
	return f.append(aiMessage)
}

func (f *FileChatMessageHistory) append(message rootSchema.BaseMessageInterface) error {
	messages, err := f.Messages()
	if err != nil {
		return err
	}

	messages = append(messages, message)
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileChatMessageHistory) Clear() error {
	err := ioutil.WriteFile(f.filePath, []byte("[]"), 0644)
	if err != nil {
		return err
	}

	return nil
}
