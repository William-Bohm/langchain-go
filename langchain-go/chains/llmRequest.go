package chains

import (
	"github.com/William-Bohm/langchain-go/langchain-go/util/requests"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var DEFAULT_HEADERS = http.Header{
	"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"},
}

type LLMRequestsChain struct {
	llm_chain         LLMChain
	requests_wrapper  requests.TextRequestsWrapper
	text_length       int
	requests_key      string
	input_key         string
	output_key        string
	arbitrary_allowed bool
	extra             string
}

func NewLLMRequestsChain() LLMRequestsChain {
	requestWrapper := requests.NewTextRequestsWrapper()
	return LLMRequestsChain{
		requests_wrapper:  requestWrapper,
		text_length:       8000,
		requests_key:      "requests_result",
		input_key:         "url",
		output_key:        "output",
		extra:             "forbid",
		arbitrary_allowed: true,
	}
}

func (chain *LLMRequestsChain) input_keys() []string {
	return []string{chain.input_key}
}

func (chain *LLMRequestsChain) output_keys() []string {
	return []string{chain.output_key}
}

func (chain *LLMRequestsChain) _call(inputs map[string]string) (map[string]string, error) {
	var other_keys map[string]interface{}
	for k, v := range inputs {
		if k != chain.input_key {
			other_keys[k] = v
		}
	}
	url := inputs[chain.input_key]
	res, err := chain.requests_wrapper.Get(url)
	if err != nil {
		return nil, err
	}
	soup := parseHTML(res)
	other_keys[chain.requests_key] = truncateString(soup, chain.text_length)
	result, err := chain.llm_chain.Predict(other_keys)
	if err != nil {
		return nil, err
	}
	return map[string]string{chain.output_key: result}, nil
}

func (chain *LLMRequestsChain) _chain_type() string {
	return "llm_requests_chain"
}

func parseHTML(htmlStr string) string {
	doc, _ := html.Parse(strings.NewReader(htmlStr))
	var f func(*html.Node)
	var sb strings.Builder
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return sb.String()
}

func truncateString(str string, num int) string {
	if len(str) > num {
		return str[0:num]
	}
	return str
}
