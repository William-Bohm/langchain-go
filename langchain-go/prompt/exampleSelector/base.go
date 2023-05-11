package exampleSelector

type BaseExampleSelector interface {
	AddExample(example map[string]string) interface{}
	SelectExamples(inputVariables map[string]string) []map[string]interface{}
	values(m map[string]interface{}) []string
}
