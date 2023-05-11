package embedding

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"log"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
	"gonum.org/v1/gonum/mat"
)

var stringToModel = map[string]openai.EmbeddingModel{
	"text-similarity-ada-001":       openai.AdaSimilarity,
	"text-similarity-babbage-001":   openai.BabbageSimilarity,
	"text-similarity-curie-001":     openai.CurieSimilarity,
	"text-similarity-davinci-001":   openai.DavinciSimilarity,
	"text-search-ada-doc-001":       openai.AdaSearchDocument,
	"text-search-ada-query-001":     openai.AdaSearchQuery,
	"text-search-babbage-doc-001":   openai.BabbageSearchDocument,
	"text-search-babbage-query-001": openai.BabbageSearchQuery,
	"text-search-curie-doc-001":     openai.CurieSearchDocument,
	"text-search-curie-query-001":   openai.CurieSearchQuery,
	"text-search-davinci-doc-001":   openai.DavinciSearchDocument,
	"text-search-davinci-query-001": openai.DavinciSearchQuery,
	"code-search-ada-code-001":      openai.AdaCodeSearchCode,
	"code-search-ada-text-001":      openai.AdaCodeSearchText,
	"code-search-babbage-code-001":  openai.BabbageCodeSearchCode,
	"code-search-babbage-text-001":  openai.BabbageCodeSearchText,
	"text-embedding-ada-002":        openai.AdaEmbeddingV2,
}

type OpenAIEmbeddingsConfig struct {
	Model              string
	Deployment         string
	EmbeddingCtxLength int
	OpenAIKey          string
	OpenAIOrganization string
	AllowedSpecial     map[string]struct{}
	DisallowedSpecial  map[string]struct{}
	ChunkSize          int
	MaxRetries         int
}

func NewOpenAIEmbeddingsConfig() *OpenAIEmbeddingsConfig {
	return &OpenAIEmbeddingsConfig{
		Model:              "text-embedding-ada-002",
		Deployment:         "text-embedding-ada-002",
		EmbeddingCtxLength: 8191,
		OpenAIKey:          "",
		OpenAIOrganization: "",
		AllowedSpecial:     map[string]struct{}{},
		DisallowedSpecial:  map[string]struct{}{"all": {}},
		ChunkSize:          1000,
		MaxRetries:         6,
	}
}

type OpenAIEmbeddings struct {
	Client *openai.Client
	Config *OpenAIEmbeddingsConfig
}

func NewOpenAIEmbeddings(config *OpenAIEmbeddingsConfig) (*OpenAIEmbeddings, error) {
	if config.OpenAIKey == "" {
		return nil, errors.New("OPENAI_API_KEY must be provided")
	}

	client := openai.NewClient(config.OpenAIKey)

	return &OpenAIEmbeddings{
		Client: client,
		Config: config,
	}, nil
}

func (oe *OpenAIEmbeddings) embedDocuments(texts []string) ([][]float64, error) {
	// TODO: Implement _get_len_safe_embeddings for large input text
	results := make([][]float64, 0, len(texts))

	g, ctx := errgroup.WithContext(context.Background())

	for i := 0; i < len(texts); i += oe.Config.ChunkSize {
		end := i + oe.Config.ChunkSize
		if end > len(texts) {
			end = len(texts)
		}

		textBatch := texts[i:end]

		g.Go(func() error {
			embeddings, err := oe.embedWithRetry(ctx, textBatch)
			if err != nil {
				return err
			}

			results = append(results, embeddings...)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}

func (oe *OpenAIEmbeddings) embedQuery(text string) ([]float64, error) {
	embedding, err := oe.embedWithRetry(context.Background(), []string{text})
	if err != nil {
		return nil, err
	}

	return embedding[0], nil
}

func (oe *OpenAIEmbeddings) embedWithRetry(ctx context.Context, texts []string) ([][]float64, error) {
	var embeddings [][]float64
	var err error

	backOff := []time.Duration{4 * time.Second, 8 * time.Second, 10 * time.Second}

	err = retry.Do(
		func() error {
			embeddings, err = oe.embed(ctx, texts)
			if err != nil {
				log.Println("embedWithRetry failed:", err)
			}
			return err
		},
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(4*time.Second),
		retry.Attempts(uint(len(backOff))),
	)

	if err != nil {
		return nil, err
	}

	return embeddings, nil
}

func (oe *OpenAIEmbeddings) embed(ctx context.Context, texts []string) ([][]float64, error) {
	cleanedTexts := make([]string, len(texts))
	for i, text := range texts {
		cleanedTexts[i] = strings.ReplaceAll(text, "\n", " ")
	}

	req := openai.EmbeddingRequest{
		Model: stringToModel[oe.Config.Model],
		Input: cleanedTexts,
	}

	resp, err := oe.Client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	embeddings := make([][]float64, len(resp.Data))
	for _, datum := range resp.Data {
		// Convert the []float32 value to []float64
		embedding := make([]float64, len(datum.Embedding))
		for i, v := range datum.Embedding {
			embedding[i] = float64(v)
		}
		embeddings = append(embeddings, embedding)
	}

	return embeddings, nil
}

func normalizeVector(vec []float64) []float64 {
	norm := mat.Norm(mat.NewVecDense(len(vec), vec), 2)
	for i := range vec {
		vec[i] /= norm
	}
	return vec
}

func main() {
	config := NewOpenAIEmbeddingsConfig()
	config.OpenAIKey = os.Getenv("OPENAI_API_KEY")

	oe, err := NewOpenAIEmbeddings(config)
	if err != nil {
		log.Fatalf("Failed to initialize OpenAIEmbeddings: %v", err)
	}

	texts := []string{
		"First test sentence.",
		"Second test sentence.",
	}

	embeddings, err := oe.embedDocuments(texts)
	if err != nil {
		log.Fatalf("Failed to embed documents: %v", err)
	}

	fmt.Println("Embeddings:")
	for i, embedding := range embeddings {
		fmt.Printf("%d: %v\n", i, normalizeVector(embedding))
	}

	query := "Test query sentence."
	queryEmbedding, err := oe.embedQuery(query)
	if err != nil {
		log.Fatalf("Failed to embed query: %v", err)
	}

	fmt.Println("Query Embedding:")
	fmt.Println(normalizeVector(queryEmbedding))
}
