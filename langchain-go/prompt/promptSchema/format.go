package promptSchema

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/util/formatting"
)

func getFormatter(template string, templateFormat string) (formatting.BaseFormatterInterface, error) {
	switch templateFormat {
	case "Text/template":
		return formatting.NewStrictFormatter(template)
	default:
		formatters := []string{
			"Text/template",
		}
		validFormats := make([]string, 0, len(formatters))
		for _, option := range formatters {
			validFormats = append(validFormats, option)
		}
		panic(fmt.Sprintf("Invalid template format. Got `%v`; should be one of %v", templateFormat, validFormats))
	}
}

func CheckValidTemplate(template string, templateFormat string, inputVariables map[string]string) error {
	// check that the template/formatter is valid
	formatter, err := getFormatter(template, templateFormat)
	if err != nil {
		return err
	}

	// Check that input variables match the schema
	err = formatter.ValidateInputVariables(inputVariables)
	if err != nil {
		return err
	}
	return nil
}

func ExecuteFormatter(templateText string, templateFormat string, inputVariables map[string]string) (string, error) {
	formatter, err := formatting.NewStrictFormatter(templateText)
	if err != nil {
		return "", err
	}
	formattedText, err := formatter.Execute(inputVariables)
	if err != nil {
		return "", err
	}

	return formattedText, nil
}
