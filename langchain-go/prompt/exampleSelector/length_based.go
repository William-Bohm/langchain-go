package exampleSelector

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptSchema" // Replace with the actual package name
)

type GetTextLengthFunc func(string) int

func GetLengthBased(text string) int {
	r := regexp.MustCompile("\n| ")
	return len(r.Split(text, -1))
}

type LengthBasedExampleSelector struct {
	Examples           []map[string]interface{}
	ExamplePrompt      promptSchema.PromptTemplate
	GetTextLength      GetTextLengthFunc
	MaxLength          int
	ExampleTextLengths []int
}

func NewLengthBasedExampleSelector(examples []map[string]interface{}, examplePrompt promptSchema.PromptTemplate) *LengthBasedExampleSelector {
	return &LengthBasedExampleSelector{
		Examples:      examples,
		ExamplePrompt: examplePrompt,
		GetTextLength: GetLengthBased,
		MaxLength:     2048,
	}
}

func (lbs *LengthBasedExampleSelector) AddExample(example map[string]interface{}) {
	lbs.Examples = append(lbs.Examples, example)
	stringExample, err := lbs.ExamplePrompt.Format(example)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lbs.ExampleTextLengths = append(lbs.ExampleTextLengths, lbs.GetTextLength(stringExample))
}

func (lbs *LengthBasedExampleSelector) CalculateExampleTextLengths() ([]int, error) {
	if len(lbs.ExampleTextLengths) > 0 {
		return lbs.ExampleTextLengths, nil
	}

	stringExamples := make([]string, len(lbs.Examples))
	for i, eg := range lbs.Examples {
		formattedExample, err := lbs.ExamplePrompt.Format(eg)
		if err != nil {
			return nil, err
		}
		stringExamples[i] = formattedExample
	}
	exampleTextLengths := make([]int, len(stringExamples))
	for i, eg := range stringExamples {
		exampleTextLengths[i] = lbs.GetTextLength(eg)
	}
	return exampleTextLengths, nil
}

func (lbs *LengthBasedExampleSelector) SelectExamples(inputVariables map[string]string) ([]map[string]interface{}, error) {
	inputs := strings.Join(values(inputVariables), " ")
	remainingLength := lbs.MaxLength - lbs.GetTextLength(inputs)
	i := 0
	examples := []map[string]interface{}{}

	exampleTextLengths, err := lbs.CalculateExampleTextLengths()
	if err != nil {
		return nil, err
	}

	for remainingLength > 0 && i < len(lbs.Examples) {
		newLength := remainingLength - exampleTextLengths[i]
		if newLength < 0 {
			break
		} else {
			examples = append(examples, lbs.Examples[i])
			remainingLength = newLength
		}
		i++
	}
	return examples, nil
}

func (lbs *LengthBasedExampleSelector) values(m map[string]interface{}) []string {
	vs := make([]string, 0, len(m))
	for _, v := range m {
		vs = append(vs, v.(string))
	}
	return vs
}
