package textSplitters

type PythonCodeTextSplitter struct {
	*RecursiveCharacterTextSplitter
}

func NewPythonCodeTextSplitter() (*PythonCodeTextSplitter, error) {
	separators := []string{
		"\nclass ",
		"\ndef ",
		"\n\tdef ",
		"\n\n",
		"\n",
		" ",
		"",
	}
	textSplitter, err := NewRecursiveCharacterTextSplitter(separators)
	if err != nil {
		return &PythonCodeTextSplitter{}, err
	}
	return &PythonCodeTextSplitter{textSplitter}, nil
}
