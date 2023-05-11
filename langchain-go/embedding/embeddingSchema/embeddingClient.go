package embeddingSchema

import (
	"bytes"
	"context"
	"errors"
	"github.com/avast/retry-go"
	"io/ioutil"
	"net/http"
	"sync"
)

type EmbeddingClient interface {
	Create(ctx context.Context) (ResponsePayload, error)
	addHeaders(req *http.Request)
}

type RequestPayload interface {
	ToJSON() ([]byte, error)
}

type ResponsePayload interface {
	FromJSON([]byte) (ResponsePayload, error)
	NewResponsePayload() ResponsePayload
}

type BaseEmbeddingClient struct {
	APIBaseURL      string
	MaxRetries      int
	client          *http.Client
	clientMutex     sync.Mutex
	ResponsePayload ResponsePayload // the responses payload type
}

// TODO: make each custom openaiClient implement the request object to handle specific authorization logic
func (c *BaseEmbeddingClient) Create(ctx context.Context, requestPayload RequestPayload) (ResponsePayload, error) {
	jsonData, err := requestPayload.ToJSON()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.APIBaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	c.AddHeaders(req)

	var resp *http.Response
	err = retry.Do(
		func() error {
			resp, err = c.getClient().Do(req)
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(uint(c.MaxRetries)),
		retry.Context(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Request failed with status: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	responsePayload := c.ResponsePayload.NewResponsePayload()
	response, err := responsePayload.FromJSON(body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *BaseEmbeddingClient) AddHeaders(req *http.Request) {
	// This method can be overridden by child structs to add custom headers
}

func NewBaseAIClient(apiBaseURL string, maxRetries int, responsePayload ResponsePayload) *BaseEmbeddingClient {
	return &BaseEmbeddingClient{
		APIBaseURL:      apiBaseURL,
		MaxRetries:      maxRetries,
		ResponsePayload: responsePayload,
	}
}

func (c *BaseEmbeddingClient) getClient() *http.Client {
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()

	if c.client == nil {
		c.client = &http.Client{}
	}

	return c.client
}
