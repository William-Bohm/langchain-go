package documentLoaders

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GitbookLoader struct {
	BaseURL         string
	WebPath         string
	LoadAllPaths    bool
	ContentSelector string
}

func NewGitbookLoader(webPage string, loadAllPaths bool, baseURL string, contentSelector string) *GitbookLoader {
	if baseURL == "" {
		baseURL = webPage
	}
	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}
	if loadAllPaths {
		webPage = baseURL + "/sitemap.xml"
	}
	return &GitbookLoader{BaseURL: baseURL, WebPath: webPage, LoadAllPaths: loadAllPaths, ContentSelector: contentSelector}
}

func (g *GitbookLoader) Load() []documentSchema.Document {
	if g.LoadAllPaths {
		soup, err := g.Scrape(g.WebPath)
		if err != nil {
			panic(err)
		}
		relativePaths := g.GetPaths(soup)
		var documents []documentSchema.Document
		for _, path := range relativePaths {
			u, err := url.Parse(path)
			if err != nil {
				panic(err)
			}
			url := g.BaseURL + u.Path
			fmt.Println("Fetching text from", url)
			soup, err := g.Scrape(url)
			if err != nil {
				panic(err)
			}
			documents = append(documents, g.GetDocument(soup, url))
		}
		return documents
	}
	soup, err := g.Scrape(g.WebPath)
	if err != nil {
		panic(err)
	}
	return []documentSchema.Document{g.GetDocument(soup, g.WebPath)}
}

func (g *GitbookLoader) Scrape(path string) (*goquery.Document, error) {
	res, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (g *GitbookLoader) GetDocument(soup *goquery.Document, customURL string) documentSchema.Document {
	pageContentRaw := soup.Find(g.ContentSelector)
	content := pageContentRaw.Text()
	title := ""
	if titleTag := pageContentRaw.Find("h1"); titleTag != nil {
		title = titleTag.Text()
	}
	metadata := map[string]interface{}{
		"source": customURL,
		"title":  title,
	}
	return documentSchema.Document{PageContent: content, Metadata: metadata}
}

func (g *GitbookLoader) GetPaths(soup *goquery.Document) []string {
	var paths []string
	soup.Find("loc").Each(func(i int, s *goquery.Selection) {
		loc := s.Text()
		paths = append(paths, loc)
	})
	return paths
}
