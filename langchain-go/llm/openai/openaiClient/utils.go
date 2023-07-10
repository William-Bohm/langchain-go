package openaiClient

func GetNumTokensForText(text string, model Model) (int, error) {
	/*
		Calculate num tokens with tiktoken-go package.
	*/
	// TODO: make a better defualt value than "text-davinci-003"

	// determine the correct encoder for the openAI model
	enc, err := GetEncodingForModel(model)
	if err != nil {
		enc, err = GetEncodingForModel("text-davinci-003")
		if err != nil {
			return 0, err
		}
	}

	// calculate the number of tokens
	ids, _, err := enc.Encode(text)
	if err != nil {
		return -1, err
	}
	return len(ids), nil
}
