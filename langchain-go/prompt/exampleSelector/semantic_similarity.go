package exampleSelector

import (
	"github.com/William-Bohm/langchain-go/langchain-go/embedding/embeddingSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/vectorstore"
	"sort"
)

type SemanticSimilarityExampleSelector struct {
	vectorstore vectorstore.VectorStore
	k           int
	exampleKeys []string
	inputKeys   []string
}

func sortedValues(values map[string]interface{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sortedVals := make([]string, len(keys))
	for i, key := range keys {
		sortedVals[i] = values[key].(string)
	}
	return sortedVals
}

func (s *SemanticSimilarityExampleSelector) addExample(example map[string]interface{}) string {
	var stringExample string
	if s.inputKeys != nil {
		filteredExample := make(map[string]interface{})
		for _, key := range s.inputKeys {
			filteredExample[key] = example[key]
		}
		stringExample = joinStrings(sortedValues(filteredExample))
	} else {
		stringExample = joinStrings(sortedValues(example))
	}
	ids, err := s.vectorstore.AddTexts([]string{stringExample}, []map[string]interface{}{example})
	if err != nil {
		return ""
	}
	return ids[0]
}

func (s *SemanticSimilarityExampleSelector) selectExamples(inputVariables map[string]interface{}) []map[string]interface{} {
	if s.inputKeys != nil {
		filteredInput := make(map[string]interface{})
		for _, key := range s.inputKeys {
			filteredInput[key] = inputVariables[key]
		}
		inputVariables = filteredInput
	}
	query := joinStrings(sortedValues(inputVariables))
	exampleDocs, err := s.vectorstore.SimilaritySearch(query, s.k)
	if err != nil {
		return []map[string]interface{}{}
	}
	examples := make([]map[string]interface{}, len(exampleDocs))
	for i, e := range exampleDocs {
		examples[i] = e.Metadata
	}
	if s.exampleKeys != nil {
		filteredExamples := make([]map[string]interface{}, len(examples))
		for i, eg := range examples {
			filteredExample := make(map[string]interface{})
			for _, key := range s.exampleKeys {
				filteredExample[key] = eg[key]
			}
			filteredExamples[i] = filteredExample
		}
		examples = filteredExamples
	}
	return examples
}

func joinStrings(strings []string) string {
	var joined string
	for _, str := range strings {
		joined += " " + str
	}
	return joined
}

func NewSemanticSimilarityExampleSelector(examples []map[string]interface{}, embeddings embeddingSchema.BaseEmbeddings, vectorstore vectorstore.VectorStore, k int, inputKeys []string) *SemanticSimilarityExampleSelector {
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
	return &SemanticSimilarityExampleSelector{
		vectorstore: vectorstore,
		k:           k,
		inputKeys:   inputKeys,
	}
}
