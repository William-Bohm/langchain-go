package formatting

import (
	"bytes"
	"fmt"
	"text/template"
)

type StrictFormatter struct {
	template *template.Template
}

func NewStrictFormatter(templateText string) (*StrictFormatter, error) {
	tmpl, err := template.New("strict").Parse(templateText)
	if err != nil {
		return nil, err
	}
	return &StrictFormatter{template: tmpl}, nil
}

func (st *StrictFormatter) Execute(inputVariables map[string]string) (string, error) {
	var buf bytes.Buffer
	err := st.template.Execute(&buf, inputVariables)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (st *StrictFormatter) ValidateTemplate(templateText string) error {
	_, err := template.New("strict").Parse(templateText)
	if err != nil {
		return err
	}

	return nil
}

func (st *StrictFormatter) ValidateInputVariables(inputVariables map[string]string) error {
	dummyInputs := make(map[string]string, len(inputVariables))
	for _, inputVariable := range inputVariables {
		dummyInputs[inputVariable] = "foo"
	}

	var buf bytes.Buffer
	err := st.template.Execute(&buf, dummyInputs)
	if err != nil {
		return fmt.Errorf("failed to validate input variables: %v", err)
	}
	return nil
}
