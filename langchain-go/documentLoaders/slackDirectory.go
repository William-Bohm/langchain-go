package documentLoaders

import (
	"archive/zip"
	"encoding/json"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type SlackDirectoryLoader struct {
	zipPath      string
	workspaceURL string
	channelIDMap map[string]string
	BaseLoader   BaseLoader
}

type Message map[string]interface{}

func NewSlackDirectoryLoader(zipPath string, workspaceURL string) *SlackDirectoryLoader {
	loader := &SlackDirectoryLoader{
		zipPath:      zipPath,
		workspaceURL: workspaceURL,
	}
	loader.channelIDMap = loader.getChannelIDMap(zipPath)
	return loader
}

func (loader *SlackDirectoryLoader) getChannelIDMap(zipPath string) map[string]string {
	r, _ := zip.OpenReader(zipPath)
	defer r.Close()
	for _, f := range r.File {
		if f.Name == "channels.json" {
			rc, _ := f.Open()
			defer rc.Close()

			bytes, _ := ioutil.ReadAll(rc)
			var channels []map[string]interface{}
			json.Unmarshal(bytes, &channels)

			channelMap := make(map[string]string)
			for _, channel := range channels {
				channelMap[channel["name"].(string)] = channel["id"].(string)
			}
			return channelMap
		}
	}
	return map[string]string{}
}

func (loader *SlackDirectoryLoader) Load() []documentSchema.Document {
	var docs []documentSchema.Document
	r, _ := zip.OpenReader(loader.zipPath)
	defer r.Close()
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".json") {
			channelName := filepath.Dir(f.Name)
			if channelName == "" {
				continue
			}
			rc, _ := f.Open()
			defer rc.Close()

			bytes, _ := ioutil.ReadAll(rc)
			var messages []Message
			json.Unmarshal(bytes, &messages)

			for _, message := range messages {
				document := loader.convertMessageToDocument(message, channelName)
				docs = append(docs, document)
			}
		}
	}
	return docs
}

func (loader *SlackDirectoryLoader) readJSON(zipFile *zip.ReadCloser, filePath string) []Message {
	for _, f := range zipFile.File {
		if f.Name == filePath {
			rc, _ := f.Open()
			defer rc.Close()

			bytes, _ := ioutil.ReadAll(rc)
			var data []Message
			json.Unmarshal(bytes, &data)
			return data
		}
	}
	return nil
}

func (loader *SlackDirectoryLoader) convertMessageToDocument(message Message, channelName string) documentSchema.Document {
	text := ""
	if val, ok := message["text"]; ok {
		text = val.(string)
	}
	metadata := loader.getMessageMetadata(message, channelName)
	return documentSchema.Document{
		PageContent: text,
		Metadata:    metadata,
	}
}

func (loader *SlackDirectoryLoader) getMessageMetadata(message Message, channelName string) map[string]interface{} {
	timestamp := ""
	if val, ok := message["ts"]; ok {
		timestamp = val.(string)
	}
	user := ""
	if val, ok := message["user"]; ok {
		user = val.(string)
	}
	source := loader.getMessageSource(channelName, user, timestamp)
	return map[string]interface{}{
		"source":    source,
		"channel":   channelName,
		"timestamp": timestamp,
		"user":      user,
	}
}

func (loader *SlackDirectoryLoader) getMessageSource(channelName string, user string, timestamp string) string {
	if loader.workspaceURL != "" {
		channelID := ""
		if val, ok := loader.channelIDMap[channelName]; ok {
			channelID = val
		}
		return loader.workspaceURL + "/archives/" + channelID + "/p" + strings.ReplaceAll(timestamp, ".", "")
	}
	return channelName + " - " + user + " - " + timestamp
}
