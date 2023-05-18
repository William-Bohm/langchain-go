package documentLoaders

import (
	"encoding/json"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type Row struct {
	SenderName  string `json:"sender_name"`
	Content     string `json:"content"`
	TimestampMs int64  `json:"timestamp_ms"`
}

func (fcl FacebookChatLoader) ConcatenateRows(row Row) string {
	sender := row.SenderName
	text := row.Content
	date := time.Unix(row.TimestampMs/1000, 0).Format("2006-01-02 15:04:05")
	return sender + " on " + date + ": " + text + "\n\n"
}

type FacebookChatLoader struct {
	FilePath string
}

func (fcl FacebookChatLoader) Load() []documentSchema.Document {
	var rawData map[string][]Row

	absPath, _ := filepath.Abs(fcl.FilePath)
	data, _ := ioutil.ReadFile(absPath)

	json.Unmarshal(data, &rawData)

	messages := rawData["messages"]
	var text strings.Builder

	for _, row := range messages {
		if row.Content != "" {
			text.WriteString(fcl.ConcatenateRows(row))
		}
	}

	metadata := map[string]interface{}{"source": absPath}

	return []documentSchema.Document{documentSchema.Document{PageContent: text.String(), Metadata: metadata}}
}

func NewFacebookChatLoader(path string) FacebookChatLoader {
	return FacebookChatLoader{FilePath: path}
}
