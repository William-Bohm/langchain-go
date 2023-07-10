package openaiClient

import (
	"fmt"
	"github.com/tiktoken-go/tokenizer"
	"strings"
)

type Model string

const DefaultModel = TextDavinci003

const (
	GPT4             Model = "gpt-4"
	GPT4_0314        Model = "gpt-4-0314"
	GPT4_32k         Model = "gpt-4-32k"
	GPT4_32k_0314    Model = "gpt-4-32k-0314"
	GPT35_Turbo      Model = "gpt-3.5-turbo"
	GPT35_Turbo_0301 Model = "gpt-3.5-turbo-0301"
	TextAda001       Model = "text-ada-001"
	Ada              Model = "ada"
	TextBabbage001   Model = "text-babbage-001"
	Babbage          Model = "babbage"
	TextCurie001     Model = "text-curie-001"
	Curie            Model = "curie"
	Davinci          Model = "davinci"
	TextDavinci003   Model = "text-davinci-003"
	TextDavinci002   Model = "text-davinci-002"
	CodeDavinci002   Model = "code-davinci-002"
	CodeDavinci001   Model = "code-davinci-001"
	CodeCushman002   Model = "code-cushman-002"
	CodeCushman001   Model = "code-cushman-001"
)

func GetEncodingForModel(model Model) (tokenizer.Codec, error) {
	/*
		helper function to determine the correct openAI token encoder depending upon the model name

		Args:
			modelname: The modelname we want to know the token encoder for

		Returns:
			The token encoder
	*/
	modelName := string(model)

	if modelName == "" {
		return tokenizer.Get(tokenizer.Cl100kBase)
	}

	switch {
	case strings.HasPrefix(string(modelName), "gpt-4"):
		return tokenizer.Get(tokenizer.Cl100kBase)
	case strings.HasPrefix(modelName, "gpt-3.5-turbo"):
		return tokenizer.Get(tokenizer.Cl100kBase)
	case strings.HasPrefix(modelName, "text-davinci-003"), strings.HasPrefix(modelName, "text-davinci-002"):
		return tokenizer.Get(tokenizer.P50kBase)
	case strings.HasPrefix(modelName, "text-davinci-001"), strings.HasPrefix(modelName, "text-curie-001"),
		strings.HasPrefix(modelName, "text-babbage-001"), strings.HasPrefix(modelName, "text-ada-001"),
		modelName == "davinci", modelName == "curie", modelName == "babbage", modelName == "ada":
		return tokenizer.Get(tokenizer.R50kBase)
	case strings.HasPrefix(modelName, "code-davinci-002"), strings.HasPrefix(modelName, "code-davinci-001"),
		strings.HasPrefix(modelName, "code-cushman-002"), strings.HasPrefix(modelName, "code-cushman-001"),
		modelName == "davinci-codex", modelName == "cushman-codex":
		return tokenizer.Get(tokenizer.P50kBase)
	case strings.HasPrefix(modelName, "text-davinci-edit-001"), strings.HasPrefix(modelName, "code-davinci-edit-001"):
		return tokenizer.Get(tokenizer.P50kEdit)
	case strings.HasPrefix(modelName, "text-embedding-ada-002"):
		return tokenizer.Get(tokenizer.Cl100kBase)
	case strings.HasPrefix(modelName, "text-similarity-davinci-001"), strings.HasPrefix(modelName, "text-similarity-curie-001"),
		strings.HasPrefix(modelName, "text-similarity-babbage-001"), strings.HasPrefix(modelName, "text-similarity-ada-001"):
		return tokenizer.Get(tokenizer.R50kBase)
	case strings.HasPrefix(modelName, "text-search-davinci-doc-001"), strings.HasPrefix(modelName, "text-search-curie-doc-001"),
		strings.HasPrefix(modelName, "text-search-babbage-doc-001"), strings.HasPrefix(modelName, "text-search-ada-doc-001"):
		return tokenizer.Get(tokenizer.R50kBase)
	case strings.HasPrefix(modelName, "code-search-babbage-code-001"), strings.HasPrefix(modelName, "code-search-ada-code-001"):
		return tokenizer.Get(tokenizer.R50kBase)
	case modelName == "gpt2":
		return tokenizer.Get(tokenizer.GPT2Enc)
	default:
		return nil, fmt.Errorf("invalid model name: %s", modelName)
	}
}

func IsValidModel(m Model) bool {
	switch m {
	case GPT4, GPT4_0314, GPT4_32k, GPT4_32k_0314, GPT35_Turbo, GPT35_Turbo_0301, TextAda001, Ada, TextBabbage001,
		Babbage, TextCurie001, Curie, Davinci, TextDavinci003, TextDavinci002, CodeDavinci002, CodeDavinci001,
		CodeCushman002, CodeCushman001:
		return true
	}
	return false
}

func (m Model) ModelNameToContextSize(modelname string) (int, error) {
	/*
		Calculate the maximum number of tokens possible to generate for a model.

		Args:
			modelname: The modelname we want to know the context size for.

		Returns:
			The maximum context size
	*/
	modelTokenMapping := map[string]int{
		"gpt-4":              8192,
		"gpt-4-0314":         8192,
		"gpt-4-32k":          32768,
		"gpt-4-32k-0314":     32768,
		"gpt-3.5-turbo":      4096,
		"gpt-3.5-turbo-0301": 4096,
		"text-ada-001":       2049,
		"ada":                2049,
		"text-babbage-001":   2040,
		"babbage":            2049,
		"text-curie-001":     2049,
		"curie":              2049,
		"davinci":            2049,
		"text-davinci-003":   4097,
		"text-davinci-002":   4097,
		"code-davinci-002":   8001,
		"code-davinci-001":   8001,
		"code-cushman-002":   2048,
		"code-cushman-001":   2048,
	}

	contextSize, ok := modelTokenMapping[modelname]
	if !ok {
		models := []string{}
		for k := range modelTokenMapping {
			models = append(models, k)
		}
		return -1, fmt.Errorf("Unknown model: %s. Please provide a valid OpenAI model name. Known models are: %s", modelname, strings.Join(models, ", "))
	}

	return contextSize, nil
}
