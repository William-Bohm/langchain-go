package documentLoaders

import (
	"encoding/json"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/ioutil"
	"net/http"
	"strings"
)

func StringifyValue(val interface{}) string {
	switch val := val.(type) {
	case string:
		return val
	case map[string]interface{}:
		return "\n" + StringifyDict(val)
	case []interface{}:
		var res []string
		for _, v := range val {
			res = append(res, StringifyValue(v))
		}
		return strings.Join(res, "\n")
	default:
		return fmt.Sprint(val)
	}
}

func StringifyDict(data map[string]interface{}) string {
	text := ""
	for key, value := range data {
		text += key + ": " + StringifyValue(value) + "\n"
	}
	return text
}

type FigmaFileLoader struct {
	AccessToken string
	Ids         string
	Key         string
}

func (f *FigmaFileLoader) ConstructFigmaAPIURL() string {
	return fmt.Sprintf("https://api.figma.com/v1/files/%s/nodes?ids=%s", f.Key, f.Ids)
}

func (f *FigmaFileLoader) GetFigmaFile() map[string]interface{} {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f.ConstructFigmaAPIURL(), nil)
	req.Header.Add("X-Figma-Token", f.AccessToken)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	return data
}

func (f *FigmaFileLoader) Load() []documentSchema.Document {
	data := f.GetFigmaFile()
	text := StringifyDict(data)
	metadata := map[string]interface{}{"source": f.ConstructFigmaAPIURL()}
	return []documentSchema.Document{{PageContent: text, Metadata: metadata}}
}

func NewFigmaFileLoader(accessToken string, ids string, key string) *FigmaFileLoader {
	return &FigmaFileLoader{AccessToken: accessToken, Ids: ids, Key: key}
}
