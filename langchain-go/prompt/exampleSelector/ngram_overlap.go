package exampleSelector

import (
	"math"

	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptSchema"
)

type NgramOverlapScoreFunc func([]string, []string) float64

func NgramOverlapScore(source []string, example []string) float64 {
	// TODO: bleu score implementation doesnt exist in golang
	return 0
}

type NGramOverlapExampleSelector struct {
	Examples      []map[string]string
	ExamplePrompt promptSchema.PromptTemplate
	Threshold     float64
}

func NewNGramOverlapExampleSelector(examples []map[string]string, examplePrompt promptSchema.PromptTemplate) *NGramOverlapExampleSelector {
	return &NGramOverlapExampleSelector{
		Examples:      examples,
		ExamplePrompt: examplePrompt,
		Threshold:     -1.0,
	}
}

func (ng *NGramOverlapExampleSelector) AddExample(example map[string]string) {
	ng.Examples = append(ng.Examples, example)
}

func (ng *NGramOverlapExampleSelector) SelectExamples(inputVariables map[string]string) ([]map[string]string, error) {
	inputs := ng.values(inputVariables)
	examples := []map[string]string{}
	k := len(ng.Examples)
	score := make([]float64, k)

	firstPromptTemplateKey := ng.ExamplePrompt.InputVariables[0]

	for i := 0; i < k; i++ {
		score[i] = NgramOverlapScore(inputs, []string{ng.Examples[i][firstPromptTemplateKey]})
	}

	for {
		argMax := argMax(score)
		if score[argMax] < ng.Threshold || math.Abs(score[argMax]-ng.Threshold) < 1e-9 {
			break
		}

		examples = append(examples, ng.Examples[argMax])
		score[argMax] = ng.Threshold - 1.0
	}

	return examples, nil
}

func (ng *NGramOverlapExampleSelector) values(m map[string]string) []string {
	vs := make([]string, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

func argMax(arr []float64) int {
	max := arr[0]
	maxIndex := 0
	for i, v := range arr {
		if v > max {
			max = v
			maxIndex = i
		}
	}
	return maxIndex
}
