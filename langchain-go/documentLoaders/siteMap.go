package documentLoaders

//
//import (
//	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
//	"regexp"
//	"strings"
//
//	"github.com/antchfx/xmlquery"
//)
//
//func defaultParsingFunction(content *xmlquery.Node) string {
//	return content.InnerText()
//}
//
//type SitemapLoader struct {
//	*WebBaseLoader
//	filterUrls      []string
//	parsingFunction func(*xmlquery.Node) string
//}
//
//func NewSitemapLoader(webPath string, filterUrls []string, parsingFunction func(*xmlquery.Node) string) *SitemapLoader {
//	if parsingFunction == nil {
//		parsingFunction = defaultParsingFunction
//	}
//
//	webLoader, _ := NewWebBaseLoader([]string{webPath}, map[string]string{})
//	return &SitemapLoader{
//		WebBaseLoader:   webLoader,
//		filterUrls:      filterUrls,
//		parsingFunction: parsingFunction,
//	}
//}
//
//func (s *SitemapLoader) ParseSitemap(soup *xmlquery.Node) []map[string]interface{} {
//	var els []map[string]interface{}
//
//	urls := xmlquery.Find(soup, "//url")
//	for _, url := range urls {
//		loc := xmlquery.FindOne(url, "loc")
//		if loc == nil {
//			continue
//		}
//
//		if s.filterUrls != nil {
//			matched := false
//			for _, r := range s.filterUrls {
//				match, _ := regexp.MatchString(r, loc.InnerText())
//				if match {
//					matched = true
//					break
//				}
//			}
//
//			if !matched {
//				continue
//			}
//		}
//
//		e := make(map[string]interface{})
//		for _, tag := range []string{"loc", "lastmod", "changefreq", "priority"} {
//			prop := xmlquery.FindOne(url, tag)
//			if prop != nil {
//				e[tag] = prop.InnerText()
//			}
//		}
//
//		els = append(els, e)
//	}
//
//	sitemaps := xmlquery.Find(soup, "//sitemap")
//	for _, sitemap := range sitemaps {
//		loc := xmlquery.FindOne(sitemap, "loc")
//		if loc == nil {
//			continue
//		}
//		soupChild, err := s.ScrapeAll([]string{loc.InnerText()}, "xml")
//		if len(soupChild) > 0 {
//			els = append(els, s.ParseSitemap(soupChild[0])...)
//		}
//	}
//
//	return els
//}
//
//func (s *SitemapLoader) Load() []documentSchema.Document {
//	soup, err := s.Scrape("xml", "")
//
//	els := s.ParseSitemap(soup)
//
//	var results []*xmlquery.Node
//	for _, el := range els {
//		if loc, ok := el["loc"]; ok {
//			resultsTemp, err := s.ScrapeAll([]string{strings.TrimSpace(loc)}, "xml")
//			results = append(results, resultsTemp[0])
//		}
//	}
//
//	var docs []documentSchema.Document
//	for i, result := range results {
//		docs = append(docs, documentSchema.Document{
//			PageContent: s.parsingFunction(result),
//			Metadata:    els[i],
//		})
//	}
//
//	return docs
//}
