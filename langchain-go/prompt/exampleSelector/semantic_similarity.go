package exampleSelector

import (
	"sort"
	"strconv"
)

type SemanticSimilarityExampleSelector struct {
	vectorstore VectorStore
	k           int
	exampleKeys []string
	inputKeys   []string
}

func sortedValues(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sortedVals := make([]string, len(keys))
	for i, key := range keys {
		sortedVals[i] = values[key]
	}
	return sortedVals
}

func (s *SemanticSimilarityExampleSelector) addExample(example map[string]string) string {
	var stringExample string
	if s.inputKeys != nil {
		filteredExample := make(map[string]string)
		for _, key := range s.inputKeys {
			filteredExample[key] = example[key]
		}
		stringExample = joinStrings(sortedValues(filteredExample))
	} else {
		stringExample = joinStrings(sortedValues(example))
	}
	ids := s.vectorstore.AddTexts([]string{stringExample}, example)
	return strconv.Itoa(ids[0])
}

func (s *SemanticSimilarityExampleSelector) selectExamples(inputVariables map[string]string) []map[string]string {
	if s.inputKeys != nil {
		filteredInput := make(map[string]string)
		for _, key := range s.inputKeys {
			filteredInput[key] = inputVariables[key]
		}
		inputVariables = filteredInput
	}
	query := joinStrings(sortedValues(inputVariables))
	exampleDocs := s.vectorstore.SimilaritySearch(query, s.k)
	examples := make([]map[string]string, len(exampleDocs))
	for i, e := range exampleDocs {
		examples[i] = e.Metadata
	}
	if s.exampleKeys != nil {
		filteredExamples := make([]map[string]string, len(examples))
		for i, eg := range examples {
			filteredExample := make(map[string]string)
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

func NewSemanticSimilarityExampleSelector(examples []map[string]string, embeddings Embeddings, vectorstoreCls VectorStore, k int, inputKeys []string) *SemanticSimilarityExampleSelector {
	var stringExamples []string
	if inputKeys != nil {
		for _, example := range examples {
			filteredExample := make(map[string]string)
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
	vectorstore := vectorstoreCls.FromTexts(stringExamples, embeddings, examples)
	return &SemanticSimilarityExampleSelector{
		vectorstore: vectorstore,
		k:           k,
		inputKeys:   inputKeys,
	}
}
