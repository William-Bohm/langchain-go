package textSplitters

import (
	"github.com/William-Bohm/langchain-go/langchain-go/llm/openai/openaiClient"
	"github.com/tiktoken-go/tokenizer"
)

type TokenTextSplitter struct {
	*BaseTextSplitter
	tokenizer         tokenizer.Codec
	allowedSpecial    interface{}
	disallowedSpecial interface{}
}

func NewTokenTextSplitter(
	modelName string,
	allowedSpecial interface{},
	disallowedSpecial interface{},
) (*TokenTextSplitter, error) {
	textSplitter, err := NewDefaultTextSplitter()
	if err != nil {
		return nil, err
	}
	model := openaiClient.Model(modelName)
	encoder, err := openaiClient.GetEncodingForModel(model)
	if err != nil {
		return nil, err
	}
	return &TokenTextSplitter{
		BaseTextSplitter:  textSplitter,
		tokenizer:         encoder,
		allowedSpecial:    allowedSpecial,
		disallowedSpecial: disallowedSpecial,
	}, nil
}

func (t *TokenTextSplitter) splitText(text string) ([]string, error) {
	var splits []string
	inputIDs, _, err := t.tokenizer.Encode(text)
	if err != nil {
		return []string{}, err
	}
	startIdx := 0
	curIdx := min(startIdx+t.chunkSize, len(inputIDs))
	chunkIDs := inputIDs[startIdx:curIdx]
	for startIdx < len(inputIDs) {
		str, err := t.tokenizer.Decode(chunkIDs)
		if err != nil {
			return []string{}, err
		}
		splits = append(splits, str)
		startIdx += t.chunkSize - t.chunkOverlap
		curIdx = min(startIdx+t.chunkSize, len(inputIDs))
		chunkIDs = inputIDs[startIdx:curIdx]
	}
	return splits, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
