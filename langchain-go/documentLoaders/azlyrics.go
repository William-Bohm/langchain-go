package documentLoaders

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
)

type AZLyricsLoader struct {
	WebBaseLoader
}

func NewAZLyricsLoader(url string) *AZLyricsLoader {
	return &AZLyricsLoader{
		WebBaseLoader: WebBaseLoader{
			URL: url,
		},
	}
}

func (l *AZLyricsLoader) Load() []*documentSchema.Document {
	soup, err := l.Scrape()
	if err != nil {
		fmt.Println("Error scraping webpage:", err)
		return nil
	}

	title := soup.Find("title").Text()
	lyrics := soup.Find("div[class='']").Eq(2).Text()
	text := title + lyrics

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
