package documentLoaders

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/ioutil"
)

type TextLoader struct {
	FilePath string
	Encoding string
}

func NewTextLoader(filePath string, encoding string) *TextLoader {
	return &TextLoader{
		FilePath: filePath,
		Encoding: encoding,
	}
}

func (loader *TextLoader) Load() []documentSchema.Document {
	content, err := ioutil.ReadFile(loader.FilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	metadata := map[string]interface{}{
		"source": loader.FilePath,
	}
	doc := documentSchema.Document{
		PageContent: string(content),
		Metadata:    metadata,
	}

	return []documentSchema.Document{doc}
}

func main() {
	loader := NewTextLoader("/path/to/text/file.txt", "utf-8")
	docs := loader.Load()

	fmt.Println("Documents:")
	for _, doc := range docs {
		fmt.Println("Page Content:", doc.PageContent)
		fmt.Println("Metadata:", doc.Metadata)
		fmt.Println()
	}
}
