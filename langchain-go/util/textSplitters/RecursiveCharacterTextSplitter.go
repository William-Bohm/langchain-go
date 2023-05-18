package textSplitters

import "strings"

type RecursiveCharacterTextSplitter struct {
	*BaseTextSplitter
	separators []string
}

func NewRecursiveCharacterTextSplitter(separators []string) (*RecursiveCharacterTextSplitter, error) {
	if len(separators) == 0 {
		separators = []string{"\n\n", "\n", " ", ""}
	}
	splitter, err := NewDefaultTextSplitter()
	if err != nil {
		return &RecursiveCharacterTextSplitter{}, err
	}
	return &RecursiveCharacterTextSplitter{
		BaseTextSplitter: splitter,
		separators:       separators,
	}, nil
}

func (r *RecursiveCharacterTextSplitter) SplitText(text string) []string {
	var finalChunks []string
	separator := r.separators[len(r.separators)-1]
	for _, s := range r.separators {
		if s == "" || strings.Contains(text, s) {
			separator = s
			break
		}
	}
	var splits []string
	if separator != "" {
		splits = strings.Split(text, separator)
	} else {
		splits = strings.Split(text, "")
	}
	var goodSplits []string
	for _, s := range splits {
		if r.lengthFunction(s) < r.chunkSize {
			goodSplits = append(goodSplits, s)
		} else {
			if len(goodSplits) > 0 {
				mergedText := r.MergeSplits(goodSplits, separator)
				finalChunks = append(finalChunks, mergedText...)
				goodSplits = []string{}
			}
			otherInfo := r.SplitText(s)
			finalChunks = append(finalChunks, otherInfo...)
		}
	}
	if len(goodSplits) > 0 {
		mergedText := r.MergeSplits(goodSplits, separator)
		finalChunks = append(finalChunks, mergedText...)
	}
	return finalChunks
}
