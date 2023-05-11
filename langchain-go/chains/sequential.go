package chains

import (
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/callbacks/callbackSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/tools"
	"strconv"
	"strings"
)

type SequentialChain struct {
	chains          []Chain
	inputVariables  []string
	outputVariables []string
	returnAll       bool
	memoryVariables []string
}

func (s *SequentialChain) InputKeys() []string {
	return s.inputVariables
}

func (s *SequentialChain) OutputKeys() []string {
	return s.outputVariables
}

func NewSequentialChain(chains []Chain, inputVariables []string, memoryVariables []string, returnAll bool) (*SequentialChain, error) {
	knownVariables := append(inputVariables, memoryVariables...)
	overlap := intersection(inputVariables, memoryVariables)
	if len(overlap) > 0 {
		return nil, errors.New("the input key(s) " + strings.Join(overlap, ",") + " are found in the Memory keys")
	}

	for _, chain := range chains {
		missingVars := difference(chain.InputKeys(), knownVariables)
		if len(missingVars) > 0 {
			return nil, errors.New("Missing required input keys: " + strings.Join(missingVars, ","))
		}
		overlap := intersection(knownVariables, chain.OutputKeys())
		if len(overlap) > 0 {
			return nil, errors.New("Chain returned keys that already exist: " + strings.Join(overlap, ","))
		}
		knownVariables = append(knownVariables, chain.OutputKeys()...)
	}

	if returnAll {
		outputVariables := difference(knownVariables, inputVariables)
		return &SequentialChain{chains, inputVariables, outputVariables, returnAll, memoryVariables}, nil
	}
	outputVariables := chains[len(chains)-1].OutputKeys()
	return &SequentialChain{chains, inputVariables, outputVariables, returnAll, memoryVariables}, nil
}

func intersection(a, b []string) []string {
	m := make(map[string]bool)
	for _, item := range a {
		m[item] = true
	}

	var result []string
	for _, item := range b {
		if _, ok := m[item]; ok {
			result = append(result, item)
		}
	}
	return result
}

func difference(a, b []string) []string {
	m := make(map[string]bool)
	for _, item := range b {
		m[item] = true
	}

	var result []string
	for _, item := range a {
		if _, ok := m[item]; !ok {
			result = append(result, item)
		}
	}
	return result
}

type SimpleSequentialChain struct {
	chains          []Chain
	stripOutputs    bool
	inputKey        string
	outputKey       string
	callbackManager callbackSchema.CallbackManager
	verbose         bool
}

func (s *SimpleSequentialChain) InputKeys() []string {
	return []string{s.inputKey}
}

func (s *SimpleSequentialChain) OutputKeys() []string {
	return []string{s.outputKey}
}

func NewSimpleSequentialChain(chains []Chain, stripOutputs bool, inputKey, outputKey string, callbackManager callbackSchema.CallbackManager, verbose bool) (*SimpleSequentialChain, error) {
	for _, chain := range chains {
		if len(chain.InputKeys()) != 1 {
			return nil, errors.New(fmt.Sprintf("Chains used in SimplePipeline should all have one input, got %v with %v inputs.", chain, len(chain.InputKeys())))
		}
		if len(chain.OutputKeys()) != 1 {
			return nil, errors.New(fmt.Sprintf("Chains used in SimplePipeline should all have one output, got %v with %v outputs.", chain, len(chain.OutputKeys())))
		}
	}
	return &SimpleSequentialChain{chains, stripOutputs, inputKey, outputKey, callbackManager, verbose}, nil
}

func (s *SimpleSequentialChain) Call(inputs map[string]string) (map[string]string, error) {
	input, ok := inputs[s.inputKey]
	if !ok {
		return nil, errors.New("Input key not found.")
	}

	items := make([]string, len(s.chains))
	for i := range s.chains {
		items[i] = strconv.Itoa(i)
	}
	colorMapping := tools.GetColorMapping(items)

	for i, chain := range s.chains {
		output, err := chain.Run(input)
		if err != nil {
			return nil, err
		}
		input = output
		if s.stripOutputs {
			input = strings.TrimSpace(input)
		}
		callbackMap := make(map[string]interface{})
		callbackMap["end"] = "\n"
		callbackMap["color"] = colorMapping[strconv.Itoa(i)]
		s.callbackManager.OnText(input, s.verbose, callbackMap)
	}
	return map[string]string{s.outputKey: input}, nil
}
