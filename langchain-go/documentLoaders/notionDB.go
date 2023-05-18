package documentLoaders

import (
	"encoding/json"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	NotionBaseURL = "https://api.notion.com/v1"
	DatabaseURL   = NotionBaseURL + "/databases/%s/query"
	PageURL       = NotionBaseURL + "/pages/%s"
	BlockURL      = NotionBaseURL + "/blocks/%s/children"
)

type NotionDBLoader struct {
	Token      string
	DatabaseID string
	Headers    map[string]string
}

func NewNotionDBLoader(integrationToken string, databaseID string) *NotionDBLoader {
	if integrationToken == "" {
		panic("integrationToken must be provided")
	}
	if databaseID == "" {
		panic("databaseID must be provided")
	}

	headers := map[string]string{
		"Authorization":  "Bearer " + integrationToken,
		"Content-Type":   "application/json",
		"Notion-Version": "2022-06-28",
	}

	return &NotionDBLoader{
		Token:      integrationToken,
		DatabaseID: databaseID,
		Headers:    headers,
	}
}

func (n *NotionDBLoader) Load() []*documentSchema.Document {
	pageIDs := n.RetrievePageIDs(map[string]interface{}{"page_size": 100})

	var documents []*documentSchema.Document
	for _, pageID := range pageIDs {
		documents = append(documents, n.LoadPage(pageID))
	}

	return documents
}

func (n *NotionDBLoader) RetrievePageIDs(queryDict map[string]interface{}) []string {
	var pages []map[string]interface{}

	for {
		data := n.Request(fmt.Sprintf(DatabaseURL, n.DatabaseID), "POST", queryDict)

		if results, ok := data["results"].([]interface{}); ok {
			for _, result := range results {
				pages = append(pages, result.(map[string]interface{}))
			}
		}

		if !data["has_more"].(bool) {
			break
		}

		queryDict["start_cursor"] = data["next_cursor"].(string)
	}

	var pageIDs []string
	for _, page := range pages {
		pageIDs = append(pageIDs, page["id"].(string))
	}

	return pageIDs
}

func (n *NotionDBLoader) LoadPage(pageID string) *documentSchema.Document {
	data := n.Request(fmt.Sprintf(PageURL, pageID), "GET", map[string]interface{}{})

	metadata := make(map[string]interface{})
	for propKey, propValueIface := range data["properties"].(map[string]interface{}) {
		propValue := propValueIface.(map[string]interface{})
		propType := propValue["type"].(string)
		var value interface{}

		switch propType {
		case "rich_text":
			value = propValue["rich_text"].([]interface{})[0].(map[string]interface{})["plain_text"]
		case "title":
			value = propValue["title"].([]interface{})[0].(map[string]interface{})["plain_text"]
		case "multi_select":
			value = []string{}
			for _, item := range propValue["multi_select"].([]interface{}) {
				value = append(value.([]string), item.(map[string]interface{})["name"].(string))
			}
		default:
			value = nil
		}

		metadata[propKey] = value
	}

	metadata["id"] = pageID

	return &documentSchema.Document{PageContent: n.LoadBlocks(pageID, 0), Metadata: metadata}
}

func (n *NotionDBLoader) LoadBlocks(blockID string, numTabs int) string {
	var resultLines []string
	curBlockID := blockID

	for curBlockID != "" {
		data := n.Request(fmt.Sprintf(BlockURL, curBlockID), "GET", map[string]interface{}{})

		for _, result := range data["results"].([]interface{}) {
			resultObj := result.(map[string]interface{})
			resultType := resultObj["type"].(string)

			if _, ok := resultObj[resultType].(map[string]interface{})["rich_text"]; !ok {
				continue
			}

			var curResultText []string
			for _, richText := range resultObj[resultType].(map[string]interface{})["rich_text"].([]interface{}) {
				if text, ok := richText.(map[string]interface{})["text"]; ok {
					curResultText = append(curResultText, strings.Repeat("\t", numTabs)+text.(map[string]interface{})["content"].(string))
				}
			}

			if resultObj["has_children"].(bool) {
				childrenText := n.LoadBlocks(resultObj["id"].(string), numTabs+1)
				curResultText = append(curResultText, childrenText)
			}

			resultLines = append(resultLines, strings.Join(curResultText, "\n"))
		}

		curBlockID = data["next_cursor"].(string)
	}

	return strings.Join(resultLines, "\n")
}

func (n *NotionDBLoader) Request(url string, method string, queryDict map[string]interface{}) map[string]interface{} {
	queryBytes, _ := json.Marshal(queryDict)
	req, _ := http.NewRequest(method, url, strings.NewReader(string(queryBytes)))
	for k, v := range n.Headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Failed to execute request with status code: %d", resp.StatusCode))
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(bodyBytes, &data)

	return data
}
