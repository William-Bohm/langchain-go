package embeddingSchema

// Embeddings is an interface for embedding models.
type BaseEmbeddings interface {
	EmbedDocuments(texts []string) ([][]float64, error)
	EmbedQuery(text string) ([]float64, error)
}
