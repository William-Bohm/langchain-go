package embedding

import (
	"fmt"
	"strings"
)

type SelfHostedEmbeddings struct {
	InferenceFn     func(pipeline string, texts []string) ([][]float64, error)
	InferenceKwargs map[string]interface{}
}

func embedDocuments(pipeline string, texts []string) ([][]float64, error) {
	// Implement the inference function to send to the remote hardware.
	return nil, nil // Replace with actual implementation
}

func NewSelfHostedEmbeddings(inferenceFn func(pipeline string, texts []string) ([][]float64, error)) *SelfHostedEmbeddings {
	return &SelfHostedEmbeddings{
		InferenceFn: inferenceFn,
	}
}

func (se *SelfHostedEmbeddings) EmbedDocuments(texts []string) ([][]float64, error) {
	cleanedTexts := make([]string, len(texts))
	for i, text := range texts {
		cleanedTexts[i] = strings.ReplaceAll(text, "\n", " ")
	}

	embeddings, err := se.InferenceFn("pipeline_reference", cleanedTexts)
	if err != nil {
		return nil, err
	}

	return embeddings, nil
}

func (se *SelfHostedEmbeddings) EmbedQuery(text string) ([]float64, error) {
	cleanedText := strings.ReplaceAll(text, "\n", " ")

	embeddings, err := se.InferenceFn("pipeline_reference", []string{cleanedText})
	if err != nil {
		return nil, err
	}

	return embeddings[0], nil
}

func main() {
	selfHostedEmbeddings := NewSelfHostedEmbeddings(embedDocuments)

	texts := []string{
		"First test sentence.",
		"Second test sentence.",
	}

	embeddings, err := selfHostedEmbeddings.EmbedDocuments(texts)
	if err != nil {
		fmt.Printf("Failed to embed documents: %v\n", err)
		return
	}

	fmt.Println("Embeddings:")
	for i, embedding := range embeddings {
		fmt.Printf("%d: %v\n", i, embedding)
	}

	query := "Test query sentence."
	queryEmbedding, err := selfHostedEmbeddings.EmbedQuery(query)
	if err != nil {
		fmt.Printf("Failed to embed query: %v\n", err)
		return
	}

	fmt.Println("Query Embedding:")
	fmt.Println(queryEmbedding)
}
