package document_compressor

import (
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/outputParser"
	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptUtils"
)

// Document represents a document object
type Document struct {
	PageContent string
}

// LLM represents a language model object
type LLM struct{}

func (l LLMChain) PredictAndParse(input map[string]interface{}) bool {
	// Implement the method to predict and parse
	// Return true or false depending on the implementation
	return true
}

// LLMChainFilter represents the main LLMChainFilter struct
type LLMChainFilter struct {
	llmChain LLMChain
	GetInput func(query string, doc Document) map[string]interface{}
}

func getDefaultChainFilterPrompt() (*promptSchema.PromptTemplate, error) {
	outputParser := outputParser.NoOutputParser{}
	outputParserMap := map[string]interface{}{"no_output_str": outputParser}
	template := promptUtils.AddInputVariablesToPrompt(outputParserMap, chainExtractPromptTemplate)

	promptTemplate, err := promptSchema.NewPromptTemplateFromTemplate(template, "default", false, outputParser, nil)
	if err != nil {
		return nil, err
	}

	return promptTemplate, nil
}

func DefaultGetInputFromChainFilter(query string, doc Document) map[string]interface{} {
	// Return the compression chain input
	input := make(map[string]interface{})
	input["question"] = query
	input["context"] = doc.PageContent
	return input
}

func (filter *LLMChainFilter) CompressDocuments(documents []Document, query string) []Document {
	filteredDocs := []Document{}
	for _, doc := range documents {
		input := filter.GetInput(query, doc)
		includeDoc := filter.llmChain.PredictAndParse(input)
		if includeDoc {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs
}

func (filter *LLMChainFilter) ACompressDocuments(documents []Document, query string) ([]Document, error) {
	// Implement the async version
	return nil, errors.New("not implemented")
}

func NewChainFilter(llm LLM, prompt string, otherParams map[string]interface{}) LLMChainFilter {
	// The function to create a new LLMChainFilter
	llmChain := LLMChain{llm: llm, prompt: prompt}
	return LLMChainFilter{llmChain: llmChain, GetInput: DefaultGetInput}
}

func main() {
	llm := LLM{}
	filter := FromLLM(llm, "default", nil)
	docs := []Document{{PageContent: "Hello, world!"}}
	query := "query"
	filtered := filter.CompressDocuments(docs, query)
	fmt.Println(filtered)
}
