package documentStore

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"github.com/gocolly/colly"
	"strings"
)

type Wikipedia struct{}

func NewWikipedia() (*Wikipedia, error) {
	// Check wikipedia package
	return &Wikipedia{}, nil
}

func (w *Wikipedia) Search(search string) (interface{}, error) {
	var result interface{}
	c := colly.NewCollector()
	c.OnHTML("#bodyContent", func(e *colly.HTMLElement) {
		result = &documentSchema.Document{
			PageContent: e.Text,
			Metadata:    map[string]interface{}{"page": e.Request.URL.String()},
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		result = "Could not find [" + search + "]. Similar: " + search
	})

	err := c.Visit("https://en.wikipedia.org/wiki/" + strings.ReplaceAll(search, " ", "_"))
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("Could not find [" + search + "]. Similar: " + search)
	}
	return result, nil
}
