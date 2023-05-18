package documentLoaders

import (
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ReadTheDocsLoader struct {
	filePath string
	encoding string
	errors   string
	bsKwargs map[string]string
}

func NewReadTheDocsLoader(filePath string, encoding string, errors string, bsKwargs map[string]string) *ReadTheDocsLoader {
	return &ReadTheDocsLoader{
		filePath: filePath,
		encoding: encoding,
		errors:   errors,
		bsKwargs: bsKwargs,
	}
}

func (r *ReadTheDocsLoader) Load() []documentSchema.Document {
	var documents []documentSchema.Document

	filepath.Walk(r.filePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		data, _ := os.ReadFile(path)

		text := r.cleanData(string(data))

		metadata := map[string]interface{}{"source": path}
		documents = append(documents, documentSchema.Document{PageContent: text, Metadata: metadata})

		return nil
	})

	return documents
}

func (r *ReadTheDocsLoader) cleanData(data string) string {
	var text strings.Builder
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(data))

	doc.Find("#main-content, [role='main']").Each(func(i int, selection *goquery.Selection) {
		for _, node := range selection.Nodes {
			text.WriteString(node.Data)
			text.WriteString("\n")
		}
	})

	return strings.Join(strings.Fields(text.String()), "\n")
}
