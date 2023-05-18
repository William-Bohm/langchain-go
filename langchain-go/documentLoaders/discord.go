package documentLoaders

import (
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"github.com/go-gota/gota/dataframe"
)

type DiscordChatLoader struct {
	BaseLoader
	ChatLog   dataframe.DataFrame
	UserIDCol string
}

func NewDiscordChatLoader(chatLog dataframe.DataFrame, userIDCol string) *DiscordChatLoader {
	return &DiscordChatLoader{ChatLog: chatLog, UserIDCol: userIDCol}
}

func (d *DiscordChatLoader) Load() []documentSchema.Document {
	var result []documentSchema.Document
	for _, row := range d.ChatLog.Records() {
		userID := ""
		metadata := map[string]interface{}{}

		for i, colName := range d.ChatLog.Names() {
			if colName == d.UserIDCol {
				userID = row[i]
			} else {
				metadata[colName] = row[i]
			}
		}
		result = append(result, documentSchema.Document{PageContent: userID, Metadata: metadata})
	}
	return result
}
