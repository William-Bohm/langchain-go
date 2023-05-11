package embedding

import (
	"errors"
	"os"
	"strings"

	"github.com/cohere-ai/cohere-go"
)

type CohereEmbeddings struct {
	Client       *cohere.Client
	Model        string
	Truncate     string
	CohereAPIKey string
}

func NewCohereEmbeddings(model string, cohereAPIKey string) (*CohereEmbeddings, error) {
	if cohereAPIKey == "" {
		cohereAPIKey = os.Getenv("COHERE_API_KEY")
	}
	if cohereAPIKey == "" {
		return nil, errors.New("COHERE_API_KEY environment variable not set")
	}

	client, err := cohere.CreateClient(cohereAPIKey)
	if err != nil {
		return nil, err
	}

	return &CohereEmbeddings{
		Client:       client,
		Model:        model,
		Truncate:     "",
		CohereAPIKey: cohereAPIKey,
	}, nil
}

func (c *CohereEmbeddings) EmbedDocuments(texts []string) ([][]float64, error) {
	embedOptions := cohere.EmbedOptions{
		Model:    c.Model,
		Texts:    texts,
		Truncate: c.Truncate,
	}
	embedResponse, err := c.Client.Embed(embedOptions)
	if err != nil {
		return nil, err
	}

	embeddings := make([][]float64, len(embedResponse.Embeddings))
	for i, e := range embedResponse.Embeddings {
		embedding := make([]float64, len(e))
		for j, v := range e {
			embedding[j] = float64(v)
		}
		embeddings[i] = embedding
	}

	return embeddings, nil
}

func (c *CohereEmbeddings) EmbedQuery(text string) ([]float64, error) {
	embeddings, err := c.EmbedDocuments([]string{text})
	if err != nil {
		return nil, err
	}
	return embeddings[0], nil
}

func TruncateFromString(truncate string) string {
	switch strings.ToUpper(truncate) {
	case "NONE":
		return ""
	case "START":
		return "start"
	case "END":
		return "end"
	default:
		return ""
	}
}
