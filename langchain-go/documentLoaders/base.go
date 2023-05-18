package documentLoaders

import (
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/util/textSplitters"
)

type BaseLoader interface {
	Load() []documentSchema.Document
	LoadAndSplit(textSplitter textSplitters.TextSplitter) []documentSchema.Document
}

type BaseLoaderImpl struct {
	BaseLoader
}

func (b *BaseLoaderImpl) Load() []documentSchema.Document {
	return []documentSchema.Document{}
}

func (b *BaseLoaderImpl) LoadAndSplit(textSplitter textSplitters.TextSplitter) ([]documentSchema.Document, error) {
	if textSplitter == nil {
		var err error
		textSplitter, err = textSplitters.NewRecursiveCharacterTextSplitter([]string{})
		if err != nil {
			return []documentSchema.Document{}, err
		}
	}
	docs := b.Load()
	return textSplitter.SplitDocuments(docs), nil
}
