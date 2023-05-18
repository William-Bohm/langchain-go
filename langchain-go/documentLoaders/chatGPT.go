package documentLoaders

import (
	"encoding/json"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"os"
	"strings"
	"time"
)

type MessagePart struct {
	Author struct {
		Role string
	}
	Content struct {
		Parts []string
	}
	CreateTime int64
}

type Mapping struct {
	Message MessagePart
}

type Data struct {
	Title   string
	Mapping map[string]Mapping
}

type ChatGPTLoader struct {
	BaseLoaderImpl
	LogFile string
	NumLogs int
}

func (c ChatGPTLoader) ConcatenateRows(message MessagePart, title string) string {
	if message.Content.Parts == nil {
		return ""
	}

	sender := "unknown"
	if message.Author.Role != "" {
		sender = message.Author.Role
	}

	text := message.Content.Parts[0]
	date := time.Unix(message.CreateTime, 0).Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s - %s on %s: %s\n\n", title, sender, date, text)
}

func (c ChatGPTLoader) Load() ([]documentSchema.Document, error) {
	file, _ := os.ReadFile(c.LogFile)
	var data []Data
	err := json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	if c.NumLogs != -1 {
		data = data[:c.NumLogs]
	}

	var documents []documentSchema.Document
	for _, d := range data {
		var parts []string
		i := 0
		for key, _ := range d.Mapping {
			if !(i == 0 && d.Mapping[key].Message.Author.Role == "system") {
				parts = append(parts, c.ConcatenateRows(d.Mapping[key].Message, d.Title))
			}
			i++
		}

		text := strings.Join(parts, "")
		metadata := map[string]interface{}{"source": c.LogFile}
		documents = append(documents, documentSchema.Document{PageContent: text, Metadata: metadata})
	}

	return documents, nil
}
