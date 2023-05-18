package documentLoaders

import (
	"errors"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type BlackboardLoader struct {
	*WebBaseLoader
	baseUrl            string
	folderPath         string
	loadAllRecursively bool
}

func NewBlackboardLoader(blackboardCourseUrl, bbrouter string, loadAllRecursively bool, basicAuth []string, cookies map[string]string) (*BlackboardLoader, error) {
	loader := &BlackboardLoader{
		loadAllRecursively: loadAllRecursively,
	}
	baseUrl, err := url.Parse(blackboardCourseUrl)
	if err != nil {
		return nil, errors.New("Invalid blackboard course url. Please provide a url that starts with https://<blackboard_url>/webapps/blackboard")
	}
	loader.baseUrl = baseUrl.Scheme + "://" + baseUrl.Host
	if cookies == nil {
		cookies = map[string]string{}
	}
	cookies["BbRouter"] = bbrouter
	loader.WebBaseLoader, err = NewWebBaseLoader(blackboardCourseUrl, basicAuth, cookies) // This function depends on your implementation
	err = checkBs4()
	if err != nil {
		return nil, err
	}
	return loader, nil
}

func checkBs4() error {
	// BeautifulSoup4 is a python library and has no equivalent in Go
	// You can remove this function or implement it depending on your needs
	return nil
}

func (loader *BlackboardLoader) load() ([]documentSchema.Document, error) {
	var documents []documentSchema.Document
	if loader.loadAllRecursively {
		soup, err := loader.Scrape() // scrape() depends on your implementation
		if err != nil {
			return nil, err
		}
		loader.folderPath, err = loader.getFolderPath(soup)
		if err != nil {
			return nil, err
		}
		relativePaths, err := loader.getPaths(soup)
		if err != nil {
			return nil, err
		}
		for _, path := range relativePaths {
			url := loader.baseUrl + path
			fmt.Printf("Fetching documents from %s\n", url)
			soup, err := loader.scrapeUrl(url)
			if err == nil {
				doc, err := loader.getDocuments(soup)
				if err == nil {
					documents = append(documents, doc...)
				}
			}
		}
	} else {
		fmt.Printf("Fetching documents from %s\n", loader.webPath)
		soup, err := loader.scrape()
		if err != nil {
			return nil, err
		}
		loader.folderPath, err = loader.getFolderPath(soup)
		if err != nil {
			return nil, err
		}
		documents, err = loader.getDocuments(soup)
		if err != nil {
			return nil, err
		}
	}
	return documents, nil
}

func (loader *BlackboardLoader) getFolderPath(soup *goquery.Document) (string, error) {
	courseName := soup.Find("span[id='crumb_1']").Text()
	if courseName == "" {
		return "", errors.New("No course name found.")
	}
	courseName = strings.TrimSpace(courseName)
	courseNameClean := strings.ReplaceAll(url.QueryEscape(courseName), "+", "_")
	folderPath := filepath.Join(".", courseNameClean)
	return folderPath, nil
}

func (loader *BlackboardLoader) getDocuments(soup *goquery.Document) ([]documentSchema.Document, error) {
	attachments, err := loader.getAttachments(soup)
	if err != nil {
		return nil, err
	}
	err = loader.downloadAttachments(attachments)
	if err != nil {
		return nil, err
	}
	return loader.loadDocuments(), nil
}

func (loader *BlackboardLoader) getAttachments(soup *goquery.Document) ([]string, error) {
	var attachments []string
	soup.Find("ul[class='contentList']").Find("ul[class='attachments']").Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && !strings.HasPrefix(href, "#") {
			attachments = append(attachments, href)
		}
	})
	if len(attachments) == 0 {
		return nil, errors.New("No content list found.")
	}
	return attachments, nil
}

func (loader *BlackboardLoader) downloadAttachments(attachments []string) error {
	os.MkdirAll(loader.folderPath, os.ModePerm)
	for _, attachment := range attachments {
		err := loader.download(attachment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (loader *BlackboardLoader) loadDocuments() []documentSchema.Document {
	loaderInstance := &DirectoryLoader{
		// parameters depends on your DirectoryLoader implementation
	}
	documents := loaderInstance.load() // load() depends on your implementation
	return documents
}

func (loader *BlackboardLoader) getPaths(soup *goquery.Document) ([]string, error) {
	var relativePaths []string
	soup.Find("ul[class='courseMenu']").Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.HasPrefix(href, "/") {
			relativePaths = append(relativePaths, href)
		}
	})
	if len(relativePaths) == 0 {
		return nil, errors.New("No course menu found.")
	}
	return relativePaths, nil
}

func (loader *BlackboardLoader) download(path string) error {
	resp, err := http.Get(loader.baseUrl + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	filename, err := loader.parseFilename(resp.Request.URL.String())
	if err != nil {
		return err
	}
	out, err := os.Create(filepath.Join(loader.folderPath, filename))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func (loader *BlackboardLoader) parseFilename(url string) (string, error) {
	urlPath := path.Base(url)
	if path.Ext(urlPath) == ".pdf" {
		return urlPath, nil
	}
	return loader.parseFilenameFromUrl(url)
}

func (loader *BlackboardLoader) parseFilenameFromUrl(url string) (string, error) {
	r := regexp.MustCompile(`filename%2A%3DUTF-8%27%27(.+)`)
	match := r.FindStringSubmatch(url)
	if len(match) == 0 {
		return "", errors.New(fmt.Sprintf("Could not parse filename from %s", url))
	}
	filename := match[1]
	if !strings.Contains(filename, ".pdf") {
		return "", errors.New(fmt.Sprintf("Incorrect file type: %s", filename))
	}
	filename = strings.Split(filename, ".pdf")[0] + ".pdf"
	filename, err := url.QueryUnescape(filename)
	if err != nil {
		return "", err
	}
	filename = strings.ReplaceAll(filename, "%20", " ")
	return filename, nil
}

func main() {
	loader := &BlackboardLoader{
		blackboardCourseUrl: "https://<YOUR BLACKBOARD URL HERE>/webapps/blackboard/content/listContent.jsp?course_id=_<YOUR COURSE ID HERE>_1&content_id=_<YOUR CONTENT ID HERE>_1&mode=reset",
		bbrouter:            "<YOUR BBROUTER COOKIE HERE>",
		loadAllRecursively:  true,
	}
	documents, err := loader.load()
	if err != nil {
		log.Fatalf("Failed to load documents: %v", err)
	}
	fmt.Printf("Loaded %d pages of PDFs from %s\n", len(documents), loader.webPath)
}
