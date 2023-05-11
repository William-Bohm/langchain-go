package chains

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/prompt/promptSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/tools"
)

type LLMChain struct {
	BaseChain
	Prompt    promptSchema.BasePromptTemplate
	LLM       llmSchema.BaseLanguageModel
	OutputKey string
}

func (c *LLMChain) InputKeys() []string {
	return c.Prompt.InputVariables
}

func (c *LLMChain) OutputKeys() []string {
	return []string{c.OutputKey}
}

func (c *LLMChain) Generate(inputList []map[string]interface{}) (llmSchema.LLMResult, error) {
	prompts, stop, err := c.prepPrompts(inputList)
	if err != nil {
		return llmSchema.LLMResult{}, err
	}
	result, err := c.LLM.Generate(prompts, stop)
	if err != nil {
		return llmSchema.LLMResult{}, err
	}

	return result, nil
}

func (c *LLMChain) prepPrompts(inputList []map[string]interface{}) ([]string, string, error) {
	var stop string
	if stopVal, ok := inputList[0]["stop"]; ok {
		stop = stopVal.(string)
	}
	var prompts []string
	for _, inputs := range inputList {
		selectedInputs := make(map[string]interface{})
		for _, k := range c.Prompt.InputVariables {
			if val, ok := inputs[k]; ok {
				selectedInputs[k] = val
			}
		}
		p, err := c.Prompt.FormatPrompt(selectedInputs)
		if err != nil {
			return nil, "", err
		}
		coloredText := tools.GetColoredText(p.ToString(), "green")
		text := "Prompt after formatting:\n" + coloredText
		callbackMap := make(map[string]interface{})
		callbackMap["verbose"] = c.verbose
		callbackMap["newline"] = "\n"
		c.callbackManager.OnText(text, callbackMap)
		if stopVal, ok := inputs["stop"]; ok && stopVal.(string) != stop {
			return nil, "", errors.New("If `stop` is present in any inputs, should be present in all.")
		}
		prompts = append(prompts, p.ToString())
	}
	return prompts, stop, nil
}

func (c *LLMChain) Apply(inputList []map[string]interface{}) ([]map[string]string, error) {
	response, err := c.Generate(inputList)
	if err != nil {
		return nil, err
	}
	return c.CreateOutputs(response), nil
}

func (c *LLMChain) CreateOutputs(response llmSchema.LLMResult) []map[string]string {
	outputs := make([]map[string]string, len(response.Generations))
	for i, generation := range response.Generations {
		outputs[i] = map[string]string{c.OutputKey: generation[0].Text}
	}
	return outputs
}

func (c *LLMChain) Predict(inputs map[string]interface{}) (string, error) {
	output, err := c.Apply([]map[string]interface{}{inputs})
	if err != nil {
		return "", err
	}
	return output[0][c.OutputKey], nil
}

func (c *LLMChain) PredictAndParse(inputs map[string]interface{}) (interface{}, error) {
	result, err := c.Predict(inputs)
	if err != nil {
		return nil, err
	}
	if c.Prompt.OutputParser != nil {
		return c.Prompt.OutputParser.Parse(result)
	} else {
		return result, nil
	}
}

func (c *LLMChain) ParseResult(result []map[string]string) ([]interface{}, error) {
	var err error
	if c.Prompt.OutputParser != nil {
		parsedResult := make([]interface{}, len(result))
		for i, res := range result {
			parsedResult[i], err = c.Prompt.OutputParser.Parse(res[c.OutputKey])
		}
		return parsedResult, err
	} else {
		//  format 'result' to fit the required []interface{} output
		var resultInterface []interface{}
		for _, m := range result {
			resultInterface = append(resultInterface, m)
		}
		return resultInterface, nil
	}
}
