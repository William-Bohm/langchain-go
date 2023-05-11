package formatting

type BaseFormatterInterface interface {
	Execute(inputVariables map[string]string) (string, error)
	ValidateTemplate(templateText string) error
	ValidateInputVariables(inputVariables map[string]string) error
}
