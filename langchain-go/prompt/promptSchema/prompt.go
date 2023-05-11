package promptSchema

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type PromptTemplate struct {
	StringPromptTemplate
	Template         string
	TemplateFormat   string //  specifie the formatter to use (currently only implementation is regex)
	ValidateTemplate bool
}

func (pt *PromptTemplate) Format(args map[string]interface{}) (string, error) {
	if pt.TemplateFormat == "Text/template" {
		tmpl, err := template.New("prompt").Parse(pt.Template)
		if err != nil {
			return "", err
		}

		var buffer strings.Builder
		err = tmpl.Execute(&buffer, args)
		if err != nil {
			return "", err
		}

		return buffer.String(), nil
	}

	return "", fmt.Errorf("unsupported template format: %s", pt.TemplateFormat)
}

func NewPromptTemplateFromExamples(examples []string, suffix string, inputVariables []string, exampleSeparator string, prefix string, additionalArgs map[string]interface{}) (*PromptTemplate, error) {
	template := strings.Join(append([]string{prefix}, append(examples, suffix)...), exampleSeparator)
	basePromptTemplate := BasePromptTemplate{}
	return &PromptTemplate{
		InputVariables:   inputVariables,
		Template:         template,
		TemplateFormat:   additionalArgs["templateFormat"].(string),
		ValidateTemplate: additionalArgs["validateTemplate"].(bool),
	}, nil
}

func NewPromptTemplateFromFile(templateFile string, inputVariables []string, templateFormat string, validateTemplate bool) (*PromptTemplate, error) {
	content, err := os.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	templateStr := string(content)
	return &PromptTemplate{
		InputVariables:   inputVariables,
		Template:         templateStr,
		TemplateFormat:   templateFormat,
		ValidateTemplate: validateTemplate,
	}, nil
}

func NewPromptTemplateFromTemplate(templateStr string, templateFormat string, validateTemplate bool, outputParser BaseOutputParser, partial map[string]interface{}) (*PromptTemplate, error) {
	// TODO: add functionality that doesnt use regex!
	if templateFormat == "" || templateFormat == "default" {
		templateFormat = "Text/template"
	}
	inputVariables := []string{}

	// Regular expression to match the variables inside curly braces
	regex := regexp.MustCompile(`\{([^}]+)\}`)
	matches := regex.FindAllStringSubmatch(templateStr, -1)

	for _, match := range matches {
		// match[1] contains the variable name inside the curly braces
		inputVariables = append(inputVariables, match[1])
	}

	stringPromptTemplate := NewStringPromptTemplate(inputVariables, outputParser, partial, "PromptTemplate")

	return &PromptTemplate{
		StringPromptTemplate: stringPromptTemplate,
		Template:             templateStr,
		TemplateFormat:       templateFormat,
		ValidateTemplate:     validateTemplate,
	}, nil
}
