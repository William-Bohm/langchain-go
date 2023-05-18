package documentLoaders

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io"
	"os"
	"strings"
)

type AirbyteJSONLoader struct {
	BaseLoaderImpl
	filePath string
}

func NewAirbyteJSONLoader(filePath string) *AirbyteJSONLoader {
	return &AirbyteJSONLoader{
		filePath: filePath,
	}
}

func (l *AirbyteJSONLoader) Load() []*documentSchema.Document {
	text := ""
	file, err := os.Open(l.filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file:", err)
			return nil
		}

		if line != "" {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(line), &data); err != nil {
				fmt.Println("Error parsing JSON:", err)
				return nil
			}

			if jsonData, ok := data["_airbyte_data"]; ok {
				text += stringifyDict(jsonData.(map[string]interface{}))
			}
		}

		if err == io.EOF {
			break
		}
	}

	metadata := map[string]interface{}{
		"source": l.filePath,
	}

	return []*documentSchema.Document{
		{
			PageContent: text,
			Metadata:    metadata,
		},
	}
}

func stringifyValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case map[string]interface{}:
		return "\n" + stringifyDict(v)
	case []interface{}:
		var strs []string
		for _, item := range v {
			strs = append(strs, stringifyValue(item))
		}
		return strings.Join(strs, "\n")
	default:
		return fmt.Sprint(v)
	}
}

func stringifyDict(data map[string]interface{}) string {
	var text strings.Builder
	for key, value := range data {
		text.WriteString(key + ": " + stringifyValue(value) + "\n")
	}
	return text.String()
}
