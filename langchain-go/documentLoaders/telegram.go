package documentLoaders

import (
	"encoding/json"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type TelegramChatLoader struct {
	filePath   string
	BaseLoader BaseLoader
}

func NewTelegramChatLoader(path string) *TelegramChatLoader {
	return &TelegramChatLoader{
		filePath: path,
	}
}

func concatenateRows(row []string) string {
	date := row[0]
	sender := row[1]
	text := row[2]
	return sender + " on " + date + ": " + text + "\n\n"
}

func (loader *TelegramChatLoader) Load() []documentSchema.Document {
	jsonFile, _ := os.Open(loader.filePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	messages := result["messages"].([]interface{})
	var records [][]string
	for _, message := range messages {
		m := message.(map[string]interface{})
		record := []string{m["date"].(string), m["text"].(string), m["from"].(string)}
		records = append(records, record)
	}

	df := dataframe.LoadRecords(records)

	df = df.Filter(dataframe.F{
		Colname:    "type",
		Comparator: series.Eq,
		Comparando: "message",
	}).Filter(dataframe.F{
		Colname:    "text",
		Comparator: series.Eq,
		Comparando: "",
	})

	df = df.Select([]string{"date", "text", "from"})
	records = df.Records()
	for i, record := range records {
		records[i] = []string{concatenateRows(record)}
	}
	df = dataframe.LoadRecords(records)
	text := strings.Join(df.Records()[0], "")

	metadata := map[string]interface{}{"source": loader.filePath}

	return []documentSchema.Document{{
		PageContent: text,
		Metadata:    metadata,
	}}
}

func main() {
	loader := NewTelegramChatLoader("/path/to/chat/dump")
	_ = loader.Load()
}
