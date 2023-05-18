package documentLoaders

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"net/http"
	"net/url"
	"strings"
)

const IFIXIT_BASE_URL = "https://www.ifixit.com/api/2.0"

type IFixitLoader struct {
	PageType string
	Id       string
	WebPath  string
}

func NewIFixitLoader(webPath string) (*IFixitLoader, error) {
	if !strings.HasPrefix(webPath, "https://www.ifixit.com") {
		return nil, errors.New("web path must start with 'https://www.ifixit.com'")
	}

	path := strings.Replace(webPath, "https://www.ifixit.com", "", -1)

	allowedPaths := []string{"/Device", "/Guide", "/Answers", "/Teardown"}

	var isValidPath bool
	for _, allowedPath := range allowedPaths {
		if strings.HasPrefix(path, allowedPath) {
			isValidPath = true
			break
		}
	}
	if !isValidPath {
		return nil, errors.New("web path must start with /Device, /Guide, /Teardown or /Answers")
	}

	pieces := strings.FieldsFunc(path, func(c rune) bool { return c == '/' })

	pageType := pieces[0]
	if pieces[0] == "Teardown" {
		pageType = "Guide"
	}

	var id string
	if pageType == "Guide" || pageType == "Answers" {
		id = pieces[2]
	} else {
		id = pieces[1]
	}

	return &IFixitLoader{PageType: pageType, Id: id, WebPath: webPath}, nil
}

func (l *IFixitLoader) Load() ([]*documentSchema.Document, error) {
	switch l.PageType {
	case "Device":
		return l.LoadDevice("", true)
	case "Guide", "Teardown":
		return l.LoadGuide("")
	case "Answers":
		return l.LoadQuestionsAndAnswers("")
	default:
		return nil, errors.New("Unknown page type: " + l.PageType)
	}
}

func LoadSuggestions(query string, docType string) ([]*documentSchema.Document, error) {
	res, err := http.Get(IFIXIT_BASE_URL + "/suggest/" + url.PathEscape(query) + "?doctypes=" + url.PathEscape(docType))
	if err != nil || res.StatusCode != 200 {
		return nil, errors.New("Could not load suggestions for " + query)
	}

	// TODO: Parse JSON response
	// data := res.json()
	// results := data["results"]

	// Here, we can't use error handling inside loop because of missing "continue" option in Go's error handling design
	output := make([]*documentSchema.Document, 0)

	// TODO: Iterate through results, fill the output
	// Note that "continue" is not needed here as Go doesn't stop execution after catching an error

	return output, nil
}

func (l *IFixitLoader) LoadDevice(urlOverride string, includeGuides bool) ([]*documentSchema.Document, error) {
	var url string
	if urlOverride == "" {
		url = IFIXIT_BASE_URL + "/wikis/CATEGORY/" + l.Id
	} else {
		url = urlOverride
	}

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	//	TODO: THIS PROBABLY DOESNT WORK, MUST DEFINE A BETTER DATA STRUCTURE TO UNMARSHAL THE JSON INTO!!!!!!!!!!
	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf("%s\n%s\n%s", data["title"], data["description"], data["contents_raw"])

	metadata := make(map[string]interface{})
	metadata["source"] = l.WebPath
	metadata["title"] = data["title"]
	documents := []*documentSchema.Document{{PageContent: text, Metadata: metadata}}

	if includeGuides {
		guideUrls := data["guides"]
		for _, guideUrl := range guideUrls.([]string) {
			loader, err := NewIFixitLoader(guideUrl)
			doc, err := loader.Load()
			if err == nil {
				documents = append(documents, doc...)
			}
		}
	}
	return documents, nil
}

func (l *IFixitLoader) LoadGuide(urlOverride string) ([]*documentSchema.Document, error) {
	var url string
	if urlOverride == "" {
		url = IFIXIT_BASE_URL + "/guides/" + l.Id
	} else {
		url = urlOverride
	}

	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return nil, errors.New("Could not load guide: " + l.WebPath)
	}

	//	TODO: THIS PROBABLY DOESNT WORK, MUST DEFINE A BETTER DATA STRUCTURE TO UNMARSHAL THE JSON INTO!!!!!!!!!!
	var data map[string]string
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	// TODO: Build the document text from the data
	docParts := make([]string, 0)
	docParts = append(docParts, "# "+data["title"], data["introduction_raw"], "\n\n###Tools Required:")

	// TODO: Handle different cases for tools and parts
	// TODO: Iterate through the steps and lines, appending to docParts
	// TODO: Append the conclusion

	// TODO: Create a new Document
	text := strings.Join(docParts, "\n")
	metadata := make(map[string]interface{})
	metadata["source"] = l.WebPath
	metadata["title"] = data["title"]

	return []*documentSchema.Document{{PageContent: text, Metadata: metadata}}, nil
}

func (l *IFixitLoader) LoadQuestionsAndAnswers(urlOverride string) ([]*documentSchema.Document, error) {
	url := l.WebPath
	if urlOverride != "" {
		url = urlOverride
	}
	loader, err := NewWebBaseLoader([]string{url}, map[string]string{})
	if err != nil {
		return nil, err
	}
	doc, err := loader.Scrape(url, "")
	if err != nil {
		return nil, err
	}

	var output []string

	title := doc.Find("h1.post-title").Text()
	output = append(output, "# "+title)

	answersHeader := doc.Find("div.post-answers-header").Text()
	if answersHeader != "" {
		output = append(output, "\n## "+answersHeader)
	}

	doc.Find(".js-answers-list .post.post-answer").Each(func(i int, s *goquery.Selection) {
		answerType := "\n### Other Answer"
		if s.Find("[itemprop='acceptedAnswer']").Length() > 0 {
			answerType = "\n### Accepted Answer"
		} else if s.HasClass("post-helpful") {
			answerType = "\n### Most Helpful Answer"
		}
		answerContent := s.Find(".post-content .post-text").Text()
		output = append(output, answerType, answerContent)
	})

	text := strings.Join(output, "\n")
	metadata := map[string]interface{}{"source": url, "title": "Questions and Answers"}

	return []*documentSchema.Document{{PageContent: text, Metadata: metadata}}, nil
}
