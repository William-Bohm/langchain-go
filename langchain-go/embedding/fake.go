package embedding

import (
	"math/rand"
)

type FakeEmbeddings struct {
	Size int
}

func NewFakeEmbeddings(size int) *FakeEmbeddings {
	return &FakeEmbeddings{Size: size}
}

func (f *FakeEmbeddings) getEmbedding() []float64 {
	embedding := make([]float64, f.Size)
	for i := range embedding {
		embedding[i] = rand.NormFloat64()
	}
	return embedding
}

func (f *FakeEmbeddings) EmbedDocuments(texts []string) [][]float64 {
	embeddings := make([][]float64, len(texts))
	for i := range texts {
		embeddings[i] = f.getEmbedding()
	}
	return embeddings
}

func (f *FakeEmbeddings) EmbedQuery(text string) []float64 {
	return f.getEmbedding()
}
