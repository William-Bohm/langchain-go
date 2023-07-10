package openai

import (
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/config/logger"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/llmSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/llm/openai/openaiClient"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
	"github.com/William-Bohm/langchain-go/langchain-go/tools/mapTools"
	"github.com/avast/retry-go"
	"os"
)

const openaiApiKeyEnvVarName = "OPENAI_API_KEY"
const openaiOrganizationEnvVarName = "OPENAI_ORGANIZATION_ID"
const openaiApiBase = "OPENAI_API_BASE"

type generatedResponse struct {
	Text         string
	FinishReason string
	LogProbs     interface{} // This can be a more specific type if you know the structure of logprobs
}

type OpenaiLLM struct {
	llmSchema.BaseLLM
	llmSchema.BaseLanguageModel
	Model              openaiClient.Model
	Client             *openaiClient.OpenAiClient
	Role               string `comment:"The role to pass to the BaseLanguageModel. ex. 'user', '"`
	ModelKwargs        map[string]interface{}
	Temperature        float64     `comment:"What sampling temperature to use."`
	MaxTokens          int         `comment:"The maximum number of tokens to generate in the completion. -1 returns as many tokens as possible given the prompt and the models maximal context size."`
	TopP               float64     `comment:"Total probability mass of tokens to consider at each step."`
	FrequencyPenalty   float64     `comment:"Penalizes repeated tokens according to frequency."`
	PresencePenalty    float64     `comment:"Penalizes repeated tokens."`
	N                  int         `comment:"How many completions to generate for each prompt."`
	BestOf             int         `comment:"Generates best_of completions server-side and returns the \"best\"."`
	OpenaiApiKey       *string     `comment:"Optional OpenAI API keys and organization."`
	OpenaiOrganization *string     `comment:"Optional OpenAI API keys and organization."`
	BatchSize          int         `comment:"Batch size to use when passing multiple documents to generate."`
	RequestTimeout     interface{} `comment:"Timeout for requests to OpenAI completion API. Default is 600 seconds."`
	LogitBias          interface{} `comment:"Adjust the probability of specific tokens being generated."`
	MaxRetries         int         `comment:"Maximum number of retries to make when generating."`
	Streaming          bool        `comment:"Whether to stream the results or not."`
	AllowedSpecial     interface{} `comment:"Set of special tokens that are allowed."`
	DisallowedSpecial  interface{} `comment:"Set of special tokens that are not allowed."`
	CompletionTokens   float64
	PromptTokens       float64
	TotalTokens        float64
}

func (o *OpenaiLLM) GetNumTokensFromMessage(messages []rootSchema.BaseMessage) (int, error) {
	var fullText string
	for _, message := range messages {
		fullText += message.Content
	}
	tokens, err := openaiClient.GetNumTokensForText(fullText, o.Model)
	if err != nil {
		return 0, err
	}
	return tokens, nil
}
func (o *OpenaiLLM) GetNumTokensFromText(text string) (int, error) {
	tokens, err := openaiClient.GetNumTokensForText(text, o.Model)
	if err != nil {
		return 0, err
	}
	return tokens, nil
}

func (o *OpenaiLLM) sendRequest(prompts []string) ([]openaiClient.CompletionResponsePayload, error) {
	// TODO: add openaiClient to openAI struct and at initialization (NewOpenaiLLM)
	// create request payload

	// send request payload to openaiClient.create

	var rawResponse []openaiClient.CompletionResponsePayload

	// Define the retry options for the createCompletion function.
	retryOpts := []retry.Option{
		retry.Attempts(uint(o.MaxRetries)),
		retry.DelayType(retry.FixedDelay),
	}

	// Wrap the createCompletion function with the retry package.
	params := o.defaultParams()
	err := retry.Do(
		func() error {
			response, err := o.Client.Create(prompts, params)
			if err != nil {
				return err
			}
			rawResponse = response
			return nil
		},
		retryOpts...,
	)

	if err != nil {
		return nil, err
	}
	return rawResponse, nil
}

// 'choices' is responses
func (o *OpenaiLLM) Generate(prompts []string, stop []string) (*llmSchema.LLMResult, error) {
	var err error
	params := o.defaultParams()
	subPrompts, err := o.GetSubPrompts(params, prompts, stop)
	if err != nil {
		return nil, err
	}
	generatedResponses := make([]generatedResponse, 0)
	tokenUsage := make(map[string]float64)

	for _, prompts := range subPrompts {
		rawResponse, err := o.sendRequest(prompts)
		for _, promptResponse := range rawResponse {
			for _, choice := range promptResponse.Choices {
				// get the text, finish reason, and log probs from the return value
				if err != nil {
					return &llmSchema.LLMResult{}, err
				}
				text := choice.Text
				finishReason := choice.FinishReason
				logProbs := choice.Logprobs

				response := generatedResponse{
					Text:         text,
					FinishReason: finishReason,
					LogProbs:     logProbs,
				}
				generatedResponses = append(generatedResponses, response)
			}

			o.updateTokenUsage(promptResponse.Usage.CompletionTokens, promptResponse.Usage.PromptTokens, promptResponse.Usage.TotalTokens) // Update token usage
		}
	}

	result, err := o.createllmSchema(generatedResponses, prompts, tokenUsage)
	if err != nil {
		return &llmSchema.LLMResult{}, err
	}

	return result, nil
}

func (o *OpenaiLLM) updateTokenUsage(completionTokens float64, promoptTokens float64, totalTokens float64) {
	o.PromptTokens = o.PromptTokens + promoptTokens
	o.CompletionTokens = o.CompletionTokens + completionTokens
	o.TotalTokens = o.TotalTokens + totalTokens
}

func (o *OpenaiLLM) createllmSchema(generatedResponses []generatedResponse, prompts []string, tokenUsage map[string]float64) (*llmSchema.LLMResult, error) {
	generations := make([][]llmSchema.Generation, len(prompts))

	for i := range prompts {
		subChoices := generatedResponses[i*o.N : (i+1)*o.N]
		generation := make([]llmSchema.Generation, len(subChoices))

		// TODO: add error handling to test if text is valid
		for j, response := range subChoices {
			text := response.Text
			finishReason := response.FinishReason
			logprobs := response.LogProbs

			generationInfo := map[string]interface{}{
				"finish_reason": finishReason,
				"logprobs":      logprobs,
			}

			generation[j] = llmSchema.Generation{
				Text:           text,
				GenerationInfo: generationInfo,
			}
		}

		generations[i] = generation
	}

	llmOutput := map[string]interface{}{
		"token_usage": tokenUsage,
		"model_name":  o.Model,
	}

	return &llmSchema.LLMResult{
		Generations: generations,
		LLMOutput:   llmOutput,
	}, nil
}

// takes a list of prompt and groups them into batches of size BatchSize
// this allows control over how many prompt are sent at a time to the API
func (o *OpenaiLLM) GetSubPrompts(params map[string]interface{}, prompts []string, stop []string) ([][]string, error) {
	var err error
	if stop != nil {
		if _, ok := params["stop"]; ok {
			return nil, errors.New("`stop` found in both the input and default params")
		}
		params["stop"] = stop
	}

	if params["max_tokens"].(int) == -1 {
		if len(prompts) != 1 {
			return nil, errors.New("max_tokens set to -1 not supported for multiple inputs")
		}
		params["max_tokens"], err = o.MaxTokensForPrompt(prompts[0]) // TODO: Change the model name to the model you are using **IMPORTANT**
		if err != nil {
			return nil, err
		}
	}

	var subPrompts [][]string
	for i := 0; i < len(prompts); i += o.BatchSize {
		end := i + o.BatchSize
		if end > len(prompts) {
			end = len(prompts)
		}
		subPrompts = append(subPrompts, prompts[i:end])
	}

	return subPrompts, nil
}

func (o *OpenaiLLM) defaultParams() map[string]interface{} {
	normalParams := map[string]interface{}{
		"temperature":       o.Temperature,
		"max_tokens":        o.MaxTokens,
		"top_p":             o.TopP,
		"frequency_penalty": o.FrequencyPenalty,
		"presence_penalty":  o.PresencePenalty,
		"n":                 o.N,
		"best_of":           o.BestOf,
		"request_timeout":   o.RequestTimeout,
		"logit_bias":        o.LogitBias,
		"streaming":         o.Streaming,
	}
	return mapTools.MergeMaps(normalParams, o.ModelKwargs)
}

func New(options ...Option) (*OpenaiLLM, error) {
	o := &OpenaiLLM{
		Model:              openaiClient.DefaultModel,
		ModelKwargs:        nil,
		Temperature:        0.7,
		MaxTokens:          256,
		TopP:               1.0,
		FrequencyPenalty:   0.0,
		PresencePenalty:    0.0,
		N:                  1,
		BestOf:             1,
		OpenaiApiKey:       nil,
		OpenaiOrganization: nil,
		RequestTimeout:     600,
		LogitBias:          nil,
		BatchSize:          20,
		MaxRetries:         2,
		Streaming:          false,
		AllowedSpecial:     nil,
		DisallowedSpecial:  "all",
	}

	for _, opt := range options {
		opt(o)
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	return o, nil
}

func NewFromMap(attrs map[string]interface{}) (*OpenaiLLM, error) {
	baseLLM, err := llmSchema.NewBaseLLM(attrs, "openai")
	if err != nil {
		logger.Error("Failed to create base BaseLanguageModel: %s", err)
		return nil, err
	}

	o := &OpenaiLLM{
		BaseLLM:            *baseLLM,
		Model:              openaiClient.DefaultModel,
		ModelKwargs:        nil,
		Temperature:        0.7,
		MaxTokens:          256,
		TopP:               1.0,
		FrequencyPenalty:   0.0,
		PresencePenalty:    0.0,
		N:                  1,
		BestOf:             1,
		OpenaiApiKey:       nil,
		OpenaiOrganization: nil,
		RequestTimeout:     600,
		LogitBias:          nil,
		BatchSize:          20,
		MaxRetries:         2,
		Streaming:          false,
		AllowedSpecial:     nil,
		DisallowedSpecial:  "all",
	}

	for key, value := range attrs {
		var opt Option

		switch key {
		case "Model":
			if val, ok := value.(string); ok {
				opt = Model(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected string, got %T", key, value)
			}
		case "Temperature":
			if val, ok := value.(float64); ok {
				opt = Temperature(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected float64, got %T", key, value)
			}
		case "MaxTokens":
			if val, ok := value.(int); ok {
				opt = MaxTokens(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected int, got %T", key, value)
			}
		case "TopP":
			if val, ok := value.(float64); ok {
				opt = TopP(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected float64, got %T", key, value)
			}
		case "FrequencyPenalty":
			if val, ok := value.(float64); ok {
				opt = FrequencyPenalty(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected float64, got %T", key, value)
			}
		case "PresencePenalty":
			if val, ok := value.(float64); ok {
				opt = PresencePenalty(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected float64, got %T", key, value)
			}
		case "N":
			if val, ok := value.(int); ok {
				opt = N(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected int, got %T", key, value)
			}
		case "BestOf":
			if val, ok := value.(int); ok {
				opt = BestOf(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected int, got %T", key, value)
			}
		case "OpenaiApiKey":
			if val, ok := value.(string); ok {
				opt = OpenaiApiKey(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected string, got %T", key, value)
			}
		case "OpenaiOrganization":
			if val, ok := value.(string); ok {
				opt = OpenaiOrganization(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected string, got %T", key, value)
			}
		case "RequestTimeout":
			opt = RequestTimeout(value)
		case "LogitBias":
			opt = LogitBias(value)
		case "BatchSize":
			if val, ok := value.(int); ok {
				opt = BatchSize(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected int, got %T", key, value)
			}
		case "MaxRetries":
			if val, ok := value.(int); ok {
				opt = MaxRetries(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected int, got %T", key, value)
			}
		case "Streaming":
			if val, ok := value.(bool); ok {
				opt = Streaming(val)
			} else {
				return nil, fmt.Errorf("invalid value type for %s: expected bool, got %T", key, value)
			}
		case "AllowedSpecial":
			opt = AllowedSpecial(value)
		case "DisallowedSpecial":
			opt = DisallowedSpecial(value)
		default:
			return nil, fmt.Errorf("unknown attribute: %s", key)
		}

		if opt != nil {
			if err := opt(o); err != nil {
				return nil, err
			}
		}
	}

	return o, nil
}

func (o *OpenaiLLM) MaxTokensForPrompt(prompt string) (int, error) {
	numTokens, err := openaiClient.GetNumTokensForText(prompt, o.Model)
	if err != nil {
		return -1, err
	}

	// get max context size for model by name
	maxContextSize, err := o.Model.ModelNameToContextSize(string(o.Model))
	if err != nil {
		return -1, err
	}

	return maxContextSize - numTokens, nil
}

/*
 * Options Pattern for openaiLLM
 *
 * This pattern provides a way to modify the default values of an openaiLLM
 * instance without needing to know all of the struct fields or create
 * a separate constructor function for each combination of field values.
 *
 * The pattern consists of defining an Option type for each field in the
 * openaiLLM struct, and a constructor function called NewOpenaiLLM that
 * takes zero or more Option values and returns a pointer to a new openaiLLM
 * instance with the specified field values.
 *
 * To use this pattern, call NewOpenaiLLM with any desired Option values to
 * create an openaiLLM instance with the desired field values. Each Option
 * value modifies a single field of the struct. If an Option value is not
 * provided, the default value for that field will be used.
 *
 * Example usage:
 *
 *   myLLM := NewOpenaiLLM(Temperature(0.5), MaxTokens(512))
 *
 * This creates a new openaiLLM instance with default values for each field,
 * except for Temperature and MaxTokens, which are set to 0.5 and 512,
 * respectively.
 */

type Option func(*OpenaiLLM) error

func Model(m string) Option {
	return func(o *OpenaiLLM) error {
		model := openaiClient.DefaultModel
		if m != "" {
			var tempModel = openaiClient.Model(m)
			if openaiClient.IsValidModel(tempModel) {
				logger.Error("Invalid model: %s", m)
				return errors.New(fmt.Sprintf("invalid model: %s", m))
			}
			model = tempModel
		}
		o.Model = model
		return nil
	}
}

func ModelKwargs(mk map[string]interface{}) Option {
	return func(o *OpenaiLLM) error {
		o.ModelKwargs = mk
		return nil
	}
}

func Temperature(t float64) Option {
	return func(o *OpenaiLLM) error {
		o.Temperature = t
		return nil
	}
}
func MaxTokens(mt int) Option {
	return func(o *OpenaiLLM) error {
		o.MaxTokens = mt
		return nil
	}
}

func TopP(tp float64) Option {
	return func(o *OpenaiLLM) error {
		o.TopP = tp
		return nil
	}
}

func FrequencyPenalty(fp float64) Option {
	return func(o *OpenaiLLM) error {
		o.FrequencyPenalty = fp
		return nil
	}
}

func PresencePenalty(pp float64) Option {
	return func(o *OpenaiLLM) error {
		o.PresencePenalty = pp
		return nil
	}
}

func N(n int) Option {
	return func(o *OpenaiLLM) error {
		o.N = n
		return nil
	}
}

func BestOf(bo int) Option {
	return func(o *OpenaiLLM) error {
		o.BestOf = bo
		return nil
	}
}

func OpenaiApiKey(key string) Option {
	return func(o *OpenaiLLM) error {
		if key == "" {
			key = os.Getenv(openaiApiKeyEnvVarName)
			if key == "" {
				return errors.New("OPENAI_API_KEY not provided or set as environment variable")
			}
		}
		o.OpenaiApiKey = &key
		return nil
	}
}

func OpenaiOrganization(org string) Option {
	return func(o *OpenaiLLM) error {
		if org == "" {
			org = os.Getenv(openaiOrganizationEnvVarName)
			if org == "" {
				return errors.New("OPENAI_ORGANIZATION not provided or set as environment variable")
			}
		}
		o.OpenaiOrganization = &org
		return nil
	}
}

func BatchSize(bs int) Option {
	return func(o *OpenaiLLM) error {
		o.BatchSize = bs
		return nil
	}
}

func RequestTimeout(rt interface{}) Option {
	return func(o *OpenaiLLM) error {
		o.RequestTimeout = rt
		return nil
	}
}

func LogitBias(lb interface{}) Option {
	return func(o *OpenaiLLM) error {
		o.LogitBias = lb
		return nil
	}
}

func MaxRetries(mr int) Option {
	return func(o *OpenaiLLM) error {
		o.MaxRetries = mr
		return nil
	}
}

func Streaming(s bool) Option {
	return func(o *OpenaiLLM) error {
		o.Streaming = s
		return nil
	}
}

func AllowedSpecial(as interface{}) Option {
	return func(o *OpenaiLLM) error {
		o.AllowedSpecial = as
		return nil
	}
}

func DisallowedSpecial(ds interface{}) Option {
	return func(o *OpenaiLLM) error {
		o.DisallowedSpecial = ds
		return nil
	}
}
