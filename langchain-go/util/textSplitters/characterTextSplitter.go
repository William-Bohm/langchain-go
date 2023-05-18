package textSplitters

import "strings"

type CharacterTextSplitter struct {
	*BaseTextSplitter
	separator string
}

func NewCharacterTextSplitter(separator string, chunkSize int, chunkOverlap int, lengthFunction func(string) int) (*CharacterTextSplitter, error) {
	textSplitter, err := NewDefaultTextSplitter()
	if err != nil {
		return nil, err
	}
	return &CharacterTextSplitter{
		BaseTextSplitter: textSplitter,
		separator:        separator,
	}, nil
}

func (c *CharacterTextSplitter) SplitText(text string) []string {
	var splits []string
	if c.separator != "" {
		splits = strings.Split(text, c.separator)
	} else {
		for _, char := range text {
			splits = append(splits, string(char))
		}
	}
	return c.MergeSplits(splits, c.separator)
}
