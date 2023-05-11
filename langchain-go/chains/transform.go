package chains

type TransformFunc func(map[string]string) (map[string]string, error)

type TransformChain struct {
	InputVariables  []string
	OutputVariables []string
	Transform       TransformFunc
}

func (t *TransformChain) InputKeys() []string {
	return t.InputVariables
}

func (t *TransformChain) OutputKeys() []string {
	return t.OutputVariables
}

func (t *TransformChain) Call(inputs map[string]string) (map[string]string, error) {
	return t.Transform(inputs)
}
