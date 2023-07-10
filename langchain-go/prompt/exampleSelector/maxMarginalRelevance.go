package exampleSelector

import (
	"github.com/William-Bohm/langchain-go/langchain-go/embedding/embeddingSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/vectorstore"
)

type MaxMarginalRelevanceExampleSelector struct {
	SemanticSimilarityExampleSelector
	fetchK int
}

func (m *MaxMarginalRelevanceExampleSelector) selectExamples(inputVariables map[string]interface{}) []map[string]interface{} {
	if m.inputKeys != nil {
		filteredInput := make(map[string]interface{})
		for _, key := range m.inputKeys {
			filteredInput[key] = inputVariables[key]
		}
		inputVariables = filteredInput
	}
	query := joinStrings(sortedValues(inputVariables))
	exampleDocs, _ := m.vectorstore.MaxMarginalRelevanceSearch(query, m.k, m.fetchK)
	examples := make([]map[string]interface{}, len(exampleDocs))
	for i, e := range exampleDocs {
		examples[i] = e.Metadata
	}
	if m.exampleKeys != nil {
		filteredExamples := make([]map[string]interface{}, len(examples))
		for i, eg := range examples {
			filteredExample := make(map[string]interface{})
			for _, key := range m.exampleKeys {
				filteredExample[key] = eg[key]
			}
			filteredExamples[i] = filteredExample
		}
		examples = filteredExamples
	}
	return examples
}

func NewMaxMarginalRelevanceExampleSelector(examples []map[string]interface{}, embeddings embeddingSchema.BaseEmbeddings, vectorstore vectorstore.VectorStore, k int, inputKeys []string, fetchK int) *MaxMarginalRelevanceExampleSelector {
	var stringExamples []string
	if inputKeys != nil {
		for _, example := range examples {
			filteredExample := make(map[string]interface{})
			for _, key := range inputKeys {
				filteredExample[key] = example[key]
			}
			stringExamples = append(stringExamples, joinStrings(sortedValues(filteredExample)))
		}
	} else {
		for _, example := range examples {
			stringExamples = append(stringExamples, joinStrings(sortedValues(example)))
		}
	}
	// TODO: properly create the vectorstore with embeddings and other values
	return &MaxMarginalRelevanceExampleSelector{
		SemanticSimilarityExampleSelector: SemanticSimilarityExampleSelector{
			vectorstore: vectorstore,
			k:           k,
			inputKeys:   inputKeys,
		},
		fetchK: fetchK,
	}
}
