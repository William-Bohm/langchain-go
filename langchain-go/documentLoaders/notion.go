package documentLoaders

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/ioutil"
	"os"
	"path/filepath"
)

type NotionDirectoryLoader struct {
	FilePath string
}

func NewNotionDirectoryLoader(path string) *NotionDirectoryLoader {
	return &NotionDirectoryLoader{
		FilePath: path,
	}
}

func (loader *NotionDirectoryLoader) Load() []documentSchema.Document {
	var docs []documentSchema.Document
	err := filepath.Walk(loader.FilePath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			metadata := map[string]interface{}{
				"source": path,
			}
			doc := documentSchema.Document{
				PageContent: string(content),
				Metadata:    metadata,
			}
			docs = append(docs, doc)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	return docs
}
