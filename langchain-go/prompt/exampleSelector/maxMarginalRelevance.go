package exampleSelector

type MaxMarginalRelevanceExampleSelector struct {
	SemanticSimilarityExampleSelector
	fetchK int
}

func (m *MaxMarginalRelevanceExampleSelector) selectExamples(inputVariables map[string]string) []map[string]string {
	if m.inputKeys != nil {
		filteredInput := make(map[string]string)
		for _, key := range m.inputKeys {
			filteredInput[key] = inputVariables[key]
		}
		inputVariables = filteredInput
	}
	query := joinStrings(sortedValues(inputVariables))
	exampleDocs := m.vectorstore.MaxMarginalRelevanceSearch(query, m.k, m.fetchK)
	examples := make([]map[string]string, len(exampleDocs))
	for i, e := range exampleDocs {
		examples[i] = e.Metadata
	}
	if m.exampleKeys != nil {
		filteredExamples := make([]map[string]string, len(examples))
		for i, eg := range examples {
			filteredExample := make(map[string]string)
			for _, key := range m.exampleKeys {
				filteredExample[key] = eg[key]
			}
			filteredExamples[i] = filteredExample
		}
		examples = filteredExamples
	}
	return examples
}

func NewMaxMarginalRelevanceExampleSelector(examples []map[string]string, embeddings Embeddings, vectorstoreClass VectorStore, k int, inputKeys []string, fetchK int) *MaxMarginalRelevanceExampleSelector {
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
	vectorstore := vectorstoreClass.FromTexts(stringExamples, embeddings, examples)
	return &MaxMarginalRelevanceExampleSelector{
		SemanticSimilarityExampleSelector: SemanticSimilarityExampleSelector{
			vectorstore: vectorstore,
			k:           k,
			inputKeys:   inputKeys,
		},
		fetchK: fetchK,
	}
}
