package openaiClient

func GetNumTokensForText(text string, model *Model) (int, error) {
	/*
		Calculate num tokens with tiktoken-go package.
	*/

	// determine the correct encoder for the openAI model
	enc, err := getEncodingForModel(*model)
	if err != nil {
		return -1, err
	}

	// calculate the number of tokens
	ids, _, err := enc.Encode(text)
	if err != nil {
		return -1, err
	}
	return len(ids), nil
}
