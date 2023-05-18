package documentLoaders

import (
	"encoding/csv"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io"
	"os"
	"strconv"
	"strings"
)

type CSVLoader struct {
	BaseLoaderImpl
	FilePath     string
	SourceColumn string
	CSVArgs      map[string]string
	Encoding     string
}

func NewCSVLoader(filePath string, sourceColumn string, csvArgs map[string]string, encoding string) *CSVLoader {
	return &CSVLoader{
		FilePath:     filePath,
		SourceColumn: sourceColumn,
		CSVArgs:      csvArgs,
		Encoding:     encoding,
	}
}

func (l *CSVLoader) Load() []*documentSchema.Document {
	file, err := os.Open(l.FilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	if l.CSVArgs != nil {
		for key, value := range l.CSVArgs {
			switch key {
			case "delimiter":
				reader.Comma = []rune(value)[0]
				// TODO: this quote funtionality may not work.
			case "quotechar":
				reader.LazyQuotes = true
			}
		}
	}

	var docs []*documentSchema.Document

	row := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading CSV file:", err)
			return nil
		}

		content := strings.Join(record, "\n")

		source := l.FilePath
		if l.SourceColumn != "" {
			val, _ := strconv.ParseInt(l.SourceColumn, 10, 64)
			source = record[val]
		}

		metadata := map[string]interface{}{
			"source": source,
			"row":    row,
		}

		doc := &documentSchema.Document{
			PageContent: content,
			Metadata:    metadata,
		}

		docs = append(docs, doc)

		row++
	}

	return docs
}
