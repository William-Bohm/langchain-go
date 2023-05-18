package documentLoaders

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var DefaultHeaderTemplate = map[string]string{
	"User-Agent":                "",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"Accept-Language":           "en-US,en;q=0.5",
	"Referer":                   "https://www.google.com/",
	"DNT":                       "1",
	"Connection":                "keep-alive",
	"Upgrade-Insecure-Requests": "1",
}

type WebBaseLoader struct {
	WebPaths          []string
	RequestsPerSecond int
	DefaultParser     string
	HeaderTemplate    map[string]string
	Client            *http.Client
}

func BuildMetadata(doc *goquery.Document, url string) map[string]interface{} {
	metadata := map[string]interface{}{"source": url}
	if title := doc.Find("title"); title != nil {
		metadata["title"] = title.Text()
	}
	if description := doc.Find("meta[name='description']"); description != nil {
		metadata["description"], _ = description.Attr("content")
	}
	if html := doc.Find("html"); html != nil {
		metadata["language"], _ = html.Attr("lang")
	}
	return metadata
}

func NewWebBaseLoader(webPaths []string, headerTemplate map[string]string) (*WebBaseLoader, error) {
	if headerTemplate == nil {
		headerTemplate = DefaultHeaderTemplate
	}

	client := &http.Client{}

	return &WebBaseLoader{
		WebPaths:          webPaths,
		RequestsPerSecond: 2,
		DefaultParser:     "html.parser",
		HeaderTemplate:    headerTemplate,
		Client:            client,
	}, nil
}

func (w *WebBaseLoader) Fetch(url string, retries int, cooldown int, backoff float64) (string, error) {
	for i := 0; i < retries; i++ {
		resp, err := w.Client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return "", err
			}
			return doc.Text(), nil
		}
		time.Sleep(time.Duration(cooldown*int(backoff)*i) * time.Second)
	}
	return "", errors.New("retry count exceeded")
}

func (w *WebBaseLoader) FetchAll(urls []string) ([]string, error) {
	sem := make(chan bool, w.RequestsPerSecond)
	var wg sync.WaitGroup
	var res []string
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sem <- true
			text, err := w.Fetch(url, 3, 2, 1.5)
			if err == nil {
				res = append(res, text)
			}
			<-sem
		}(url)
	}
	wg.Wait()
	return res, nil
}

func (w *WebBaseLoader) ScrapeAll(urls []string, parser string) ([]*goquery.Document, error) {
	texts, err := w.FetchAll(urls)
	if err != nil {
		return nil, err
	}
	var docs []*goquery.Document
	for _, text := range texts {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (w *WebBaseLoader) Scrape(url string, parser string) (*goquery.Document, error) {
	text, err := w.Fetch(url, 3, 2, 1.5)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	return doc, err
}

func (w *WebBaseLoader) Load() ([]documentSchema.Document, error) {
	var docs []documentSchema.Document
	for _, path := range w.WebPaths {
		doc, err := w.Scrape(path, w.DefaultParser)
		if err != nil {
			return nil, err
		}
		text := doc.Text()
		metadata := BuildMetadata(doc, path)
		docs = append(docs, documentSchema.Document{PageContent: text, Metadata: metadata})
	}
	return docs, nil
}

func (w *WebBaseLoader) ALoad() ([]documentSchema.Document, error) {
	soups, err := w.ScrapeAll(w.WebPaths, w.DefaultParser)
	if err != nil {
		return nil, err
	}
	var docs []documentSchema.Document
	for i, soup := range soups {
		text := soup.Text()
		metadata := BuildMetadata(soup, w.WebPaths[i])
		docs = append(docs, documentSchema.Document{PageContent: text, Metadata: metadata})
	}
	return docs, nil
}
