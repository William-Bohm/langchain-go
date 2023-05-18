package documentLoaders

import (
	"encoding/csv"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io"
	"os"
	"strings"
)

type CoNLLULoader struct {
	BaseLoaderImpl
	FilePath string
}

func NewCoNLLULoader(filePath string) *CoNLLULoader {
	return &CoNLLULoader{
		FilePath: filePath,
	}
}

func (l *CoNLLULoader) Load() []*documentSchema.Document {
	file, err := os.Open(l.FilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	var lines [][]string
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading file:", err)
			return nil
		}

		if len(line) > 1 {
			lines = append(lines, line)
		}
	}

	var text strings.Builder
	for i, line := range lines {
		if line[9] == "SpaceAfter=No" || i == len(lines)-1 {
			text.WriteString(line[1])
		} else {
			text.WriteString(line[1] + " ")
		}
	}

	metadata := map[string]interface{}{
		"source": l.FilePath,
	}

	return []*documentSchema.Document{
		{
			PageContent: text.String(),
			Metadata:    metadata,
		},
	}
}
