package textSplitters

type SpacyTextSplitter struct {
	*BaseTextSplitter
	tokenizer func(string) []string
	separator string
}

func NewSpacyTextSplitter(separator string, pipeline string, chunkSize int, chunkOverlap int, lengthFunction func(string) int) (*SpacyTextSplitter, error) {
	textSplitter, err := NewDefaultTextSplitter()
	if err != nil {
		return &SpacyTextSplitter{}, err
	}
	return &SpacyTextSplitter{
		BaseTextSplitter: textSplitter,
		tokenizer:        SpacySentTokenize,
		separator:        separator,
	}, nil
}

func (s *SpacyTextSplitter) SplitText(text string) []string {
	splits := s.tokenizer(text)
	return s.MergeSplits(splits, s.separator)
}
