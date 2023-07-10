package document_compressor

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/outputParser"
	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptUtils"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
)

type LLMChain struct {
	// Define the fields of the LLMChain struct as needed.
}

type LLMChainExtractor struct {
	LLMChain LLMChain
	GetInput func(query string, doc rootSchema.Document) (map[string]interface{}, error)
}

func DefaultGetInputFromChainExtract(query string, doc rootSchema.Document) map[string]interface{} {
	return map[string]interface{}{
		"question": query,
		"context":  doc.PageContent,
	}
}

func getDefaultChainExtractPrompt() (*promptSchema.PromptTemplate, error) {
	outputParser := outputParser.NoOutputParser{}
	outputParserMap := map[string]interface{}{"no_output_str": outputParser}
	template := promptUtils.AddInputVariablesToPrompt(outputParserMap, chainExtractPromptTemplate)

	promptTemplate, err := promptSchema.NewPromptTemplateFromTemplate(template, "default", false, &outputParser, nil)
	if err != nil {
		return nil, err
	}

	return promptTemplate, nil
}

func (e *LLMChainExtractor) CompressDocuments(documents []rootSchema.Document, query string) ([]rootSchema.Document, error) {
	compressedDocs := []rootSchema.Document{}
	for _, doc := range documents {
		input, err := e.GetInput(query, doc)
		if err != nil {
			return nil, err
		}
		output, err := e.LLMChain.PredictAndParse(input)
		if err != nil {
			return nil, err
		}
		if len(output) == 0 {
			continue
		}
		compressedDocs = append(compressedDocs, rootSchema.Document{PageContent: output, Metadata: doc.Metadata})
	}
	return compressedDocs, nil
}

func (e *LLMChainExtractor) ACompressDocuments(documents []rootSchema.Document, query string) ([]rootSchema.Document, error) {
	return nil, errors.New("not implemented")
}

type BaseLanguageModel struct {
	// Define the fields of the BaseLanguageModel struct as needed.
}

func NewChainExtractor(
	llm BaseLanguageModel,
	prompt *promptSchema.PromptTemplate,
	getInput func(query string, doc rootSchema.Document) (map[string]interface{}, error),
) *LLMChainExtractor {
	if prompt == nil {
		noOutputParser := outputParser.NoOutputParser{NoOutputStr: "NO_OUTPUT"}
		prompt = &PromptTemplate{
			Template:       "some_template", // Replace with the actual template
			InputVariables: []string{"question", "context"},
			OutputParser:   noOutputParser,
		}
	}
	if getInput == nil {
		getInput = DefaultGetInputFromChainExtract()
	}
	llmChain := LLMChain{} // Initialize the LLMChain with the required fields
	return &LLMChainExtractor{LLMChain: llmChain, GetInput: getInput}
}
