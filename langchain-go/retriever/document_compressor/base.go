package document_compressor

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
)

type BaseDocumentCompressor interface {
	CompressDocuments(documents []rootSchema.Document, query string) ([]rootSchema.Document, error)
}

type BaseDocumentTransformer interface {
	TransformDocuments(documents []rootSchema.Document) ([]rootSchema.Document, error)
}

type DocumentCompressorPipeline struct {
	Transformers []interface{}
}

func (p *DocumentCompressorPipeline) CompressDocuments(documents []rootSchema.Document, query string) ([]rootSchema.Document, error) {
	for _, transformer := range p.Transformers {
		switch t := transformer.(type) {
		case BaseDocumentCompressor:
			var err error
			documents, err = t.CompressDocuments(documents, query)
			if err != nil {
				return nil, err
			}
		case BaseDocumentTransformer:
			var err error
			documents, err = t.TransformDocuments(documents)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("got unexpected transformer type")
		}
	}
	return documents, nil
}
