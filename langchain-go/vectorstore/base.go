package vectorstore

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
)

import "context"

type VectorStore interface {
	AddTexts([]string, []map[string]interface{}) ([]string, error)
	AddDocuments([]documentSchema.Document) ([]string, error)
	SimilaritySearch(string, int) ([]documentSchema.Document, error)
	SimilaritySearchWithRelevanceScores(string, int) ([]documentSchema.Document, float64, error)
	SimilaritySearchByVector([]float64, int) ([]documentSchema.Document, error)
	MaxMarginalRelevanceSearch(string, int, int) ([]documentSchema.Document, error)
	MaxMarginalRelevanceSearchByVector([]float64, int, int) ([]documentSchema.Document, error)
	FromDocuments([]documentSchema.Document) (VectorStore, error)
	FromTexts([]string) (VectorStore, error)
	AsRetriever() (VectorStoreRetriever, error)
}

type VectorStoreRetriever struct {
	VectorStore  VectorStore
	SearchType   string
	SearchKwargs map[string]interface{}
}

func NewVectorStoreRetriever(vs VectorStore, st string, sk map[string]interface{}) (*VectorStoreRetriever, error) {
	if st != "similarity" && st != "mmr" {
		return nil, errors.New("search_type of " + st + " not allowed.")
	}

	vsr := &VectorStoreRetriever{
		VectorStore:  vs,
		SearchType:   st,
		SearchKwargs: sk,
	}
	return vsr, nil
}

func (vsr *VectorStoreRetriever) GetRelevantDocuments(ctx context.Context, query string) ([]documentSchema.Document, error) {
	var docs []documentSchema.Document
	var err error

	if vsr.SearchType == "similarity" {
		docs, err = vsr.VectorStore.SimilaritySearch(query, 0)
	} else if vsr.SearchType == "mmr" {
		docs, err = vsr.VectorStore.MaxMarginalRelevanceSearch(query, 0, 0)
	} else {
		return nil, errors.New("search_type of " + vsr.SearchType + " not allowed.")
	}

	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (vsr *VectorStoreRetriever) AddDocuments(ctx context.Context, docs []documentSchema.Document) ([]string, error) {
	return vsr.VectorStore.AddDocuments(docs)
}
