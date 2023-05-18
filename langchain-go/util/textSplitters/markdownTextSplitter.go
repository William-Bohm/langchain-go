package textSplitters

type MarkdownTextSplitter struct {
	*RecursiveCharacterTextSplitter
}

func NewMarkdownTextSplitter() (*MarkdownTextSplitter, error) {
	separators := []string{
		"\n## ",
		"\n### ",
		"\n#### ",
		"\n##### ",
		"\n###### ",
		"```\n\n",
		"\n\n***\n\n",
		"\n\n---\n\n",
		"\n\n___\n\n",
		"\n\n",
		"\n",
		" ",
		"",
	}
	textSplitter, err := NewRecursiveCharacterTextSplitter(separators)
	if err != nil {
		return &MarkdownTextSplitter{}, err
	}
	return &MarkdownTextSplitter{textSplitter}, nil
}
