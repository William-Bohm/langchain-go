package documentLoaders

import (
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/ioutil"
	"net/http"
	"strings"
)

type GutenbergLoader struct {
	FilePath string
}

func NewGutenbergLoader(filePath string) *GutenbergLoader {
	if !strings.HasPrefix(filePath, "https://www.gutenberg.org") {
		panic("file path must start with 'https://www.gutenberg.org'")
	}

	if !strings.HasSuffix(filePath, ".txt") {
		panic("file path must end with '.txt'")
	}

	return &GutenbergLoader{
		FilePath: filePath,
	}
}

func (l *GutenbergLoader) Load() []*documentSchema.Document {
	resp, err := http.Get(l.FilePath)
	if err != nil {
		panic("failed to load file: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("failed to read response body: " + err.Error())
	}

	text := string(body)
	metadata := map[string]interface{}{
		"source": l.FilePath,
	}

	return []*documentSchema.Document{
		{
			PageContent: text,
			Metadata:    metadata,
		},
	}
}
