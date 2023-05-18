package textSplitters

type LatexTextSplitter struct {
	*RecursiveCharacterTextSplitter
}

func NewLatexTextSplitter() (*LatexTextSplitter, error) {
	separators := []string{
		"\n\\chapter{",
		"\n\\section{",
		"\n\\subsection{",
		"\n\\subsubsection{",
		"\n\\begin{enumerate}",
		"\n\\begin{itemize}",
		"\n\\begin{description}",
		"\n\\begin{list}",
		"\n\\begin{quote}",
		"\n\\begin{quotation}",
		"\n\\begin{verse}",
		"\n\\begin{verbatim}",
		"\n\\begin{align}",
		"$$",
		"$",
		" ",
		"",
	}
	textSplitter, err := NewRecursiveCharacterTextSplitter(separators)
	if err != nil {
		return &LatexTextSplitter{}, err
	}
	return &LatexTextSplitter{textSplitter}, nil
}
