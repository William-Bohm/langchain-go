package textSplitters

type NLTKTextSplitter struct {
	*BaseTextSplitter
	tokenizer func(string) []string
	separator string
}

func NewNLTKTextSplitter(separator string, chunkSize int, chunkOverlap int, lengthFunction func(string) int) (*NLTKTextSplitter, error) {
	textSplitter, err := NewDefaultTextSplitter()
	if err != nil {
		return &NLTKTextSplitter{}, err
	}
	return &NLTKTextSplitter{
		BaseTextSplitter: textSplitter,
		tokenizer:        SentTokenize,
		separator:        separator,
	}, nil
}

func (n *NLTKTextSplitter) SplitText(text string) []string {
	splits := n.tokenizer(text)
	return n.MergeSplits(splits, n.separator)
}
