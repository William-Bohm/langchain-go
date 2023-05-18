package documentLoaders

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"strings"
)

type CollegeConfidentialLoader struct {
	WebBaseLoader
}

func NewCollegeConfidentialLoader(url string) *CollegeConfidentialLoader {
	return &CollegeConfidentialLoader{
		WebBaseLoader: WebBaseLoader{
			URL: url,
		},
	}
}

func (l *CollegeConfidentialLoader) Load() []*documentSchema.Document {
	soup, err := l.Scrape()
	if err != nil {
		fmt.Println("Error scraping webpage:", err)
		return nil
	}

	text := soup.Find("main[class='skin-handler']").Text()
	text = strings.TrimSpace(text)

	metadata := map[string]interface{}{
		"source": l.WebPath,
	}

	return []*documentSchema.Document{
		{
			PageContent: text,
			Metadata:    metadata,
		},
	}
}
