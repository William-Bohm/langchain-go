package textSplitters

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"strings"
)

type TextSplitter interface {
	SplitText(text string) []string
	CreateDocuments(texts []string, metadatas []map[string]interface{}) []documentSchema.Document
	SplitDocuments(documents []documentSchema.Document) []documentSchema.Document
	JoinDocs(docs []string, separator string) string
	MergeSplits(splits []string, separator string) []string
	TransformDocuments(documents []documentSchema.Document) []documentSchema.Document
}

type BaseTextSplitter struct {
	TextSplitter
	chunkSize      int
	chunkOverlap   int
	lengthFunction func(string) int
}

func NewTextSplitter(chunkSize int, chunkOverlap int, lengthFunction func(string) int) (*BaseTextSplitter, error) {
	if chunkOverlap > chunkSize {
		return nil, errors.New("Got a larger chunk overlap (" + string(chunkOverlap) + ") than chunk size (" + string(chunkSize) + "), should be smaller.")
	}
	return &BaseTextSplitter{
		chunkSize:      chunkSize,
		chunkOverlap:   chunkOverlap,
		lengthFunction: lengthFunction,
	}, nil
}

func NewDefaultTextSplitter() (*BaseTextSplitter, error) {
	return &BaseTextSplitter{
		chunkSize:      4000,
		chunkOverlap:   200,
		lengthFunction: len,
	}, nil
}

type BaseDocumentTransformer interface {
	SplitText(text string) []string
	CreateDocuments(texts []string, metadatas []map[string]interface{}) []documentSchema.Document
	SplitDocuments(documents []documentSchema.Document) []documentSchema.Document
}

func (t *BaseTextSplitter) SplitText(text string) []string {
	// Implemented in child classes
	return nil
}

func (t *BaseTextSplitter) CreateDocuments(texts []string, metadatas []map[string]interface{}) []documentSchema.Document {
	if metadatas == nil {
		metadatas = make([]map[string]interface{}, len(texts))
	}

	var documents []documentSchema.Document
	for i, text := range texts {
		for _, chunk := range t.SplitText(text) {
			newDoc := NewDocument(chunk, metadatas[i])
			documents = append(documents, newDoc)
		}
	}
	return documents
}

func (t *BaseTextSplitter) SplitDocuments(documents []documentSchema.Document) []documentSchema.Document {
	var texts []string
	var metadatas []map[string]interface{}
	for _, doc := range documents {
		texts = append(texts, doc.PageContent)
		metadatas = append(metadatas, doc.Metadata)
	}
	return t.CreateDocuments(texts, metadatas)
}

func (t *BaseTextSplitter) JoinDocs(docs []string, separator string) string {
	text := strings.Join(docs, separator)
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	} else {
		return text
	}
}

func (t *BaseTextSplitter) MergeSplits(splits []string, separator string) []string {
	separatorLen := t.lengthFunction(separator)

	var docs []string
	var currentDoc []string
	total := 0
	for _, d := range splits {
		_len := t.lengthFunction(d)
		if total+_len+(separatorLen*len(currentDoc)) > t.chunkSize {
			if total > t.chunkSize {
				// log warning
			}
			if len(currentDoc) > 0 {
				doc := t.JoinDocs(currentDoc, separator)
				if doc != "" {
					docs = append(docs, doc)
				}
				for total > t.chunkOverlap || (total+_len+(separatorLen*len(currentDoc)) > t.chunkSize && total > 0) {
					total -= t.lengthFunction(currentDoc[0]) + (separatorLen * (len(currentDoc) - 1))
					currentDoc = currentDoc[1:]
				}
			}
			currentDoc = append(currentDoc, d)
			total += _len + (separatorLen * (len(currentDoc) - 1))
		}
		doc := t.JoinDocs(currentDoc, separator)
		if doc != "" {
			docs = append(docs, doc)
		}
	}
	return docs
}

func (t *BaseTextSplitter) TransformDocuments(documents []documentSchema.Document) []documentSchema.Document {
	return t.SplitDocuments(documents)
}

func NewDocument(pageContent string, metadata map[string]interface{}) documentSchema.Document {
	return documentSchema.Document{PageContent: pageContent, Metadata: metadata}
}
