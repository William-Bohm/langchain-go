package retriever

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/valyala/fasthttp"
)

type Document struct {
	PageContent string                 `json:"page_content"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type ChatGPTPluginRetriever struct {
	URL          string
	BearerToken  string
	TopK         int
	Filter       map[string]interface{}
	AsyncSession *fasthttp.Client
}

func (r *ChatGPTPluginRetriever) createRequest(query string) (string, map[string]interface{}, map[string]string) {
	url := fmt.Sprintf("%s/query", r.URL)
	jsonBody := map[string]interface{}{
		"queries": []map[string]interface{}{
			{
				"query":  query,
				"filter": r.Filter,
				"top_k":  r.TopK,
			},
		},
	}
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", r.BearerToken),
	}
	return url, jsonBody, headers
}

func (r *ChatGPTPluginRetriever) GetRelevantDocuments(query string) ([]Document, error) {
	url, jsonBody, headers := r.createRequest(query)

	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	results := data["results"].([]interface{})[0].(map[string]interface{})["results"].([]interface{})
	var docs []Document

	for _, d := range results {
		docData := d.(map[string]interface{})
		content := docData["text"].(string)
		delete(docData, "text")
		docs = append(docs, Document{PageContent: content, Metadata: docData})
	}

	return docs, nil
}
