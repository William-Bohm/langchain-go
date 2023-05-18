package documentLoaders

import (
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ObsidianLoader struct {
	filePath        string
	encoding        string
	collectMetadata bool
}

var FrontMatterRegex = regexp.MustCompile(`(?ms)^---\n(.*?)\n---\n`)

func NewObsidianLoader(path string, encoding string, collectMetadata bool) ObsidianLoader {
	return ObsidianLoader{filePath: path, encoding: encoding, collectMetadata: collectMetadata}
}

func (l *ObsidianLoader) parseFrontMatter(content string) map[string]interface{} {
	frontMatter := make(map[string]interface{})

	if !l.collectMetadata {
		return frontMatter
	}

	match := FrontMatterRegex.FindStringSubmatch(content)
	if match != nil {
		lines := strings.Split(match[1], "\n")
		for _, line := range lines {
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
				frontMatter[key] = value
			}
		}
	}

	return frontMatter
}

func (l *ObsidianLoader) removeFrontMatter(content string) string {
	if !l.collectMetadata {
		return content
	}

	return FrontMatterRegex.ReplaceAllString(content, "")
}

func (l *ObsidianLoader) load() []documentSchema.Document {
	var docs []documentSchema.Document
	files, _ := filepath.Glob(filepath.Join(l.filePath, "**/*.md"))

	for _, file := range files {
		data, _ := os.ReadFile(file)
		text := string(data)

		frontMatter := l.parseFrontMatter(text)
		text = l.removeFrontMatter(text)

		metadata := frontMatter
		metadata["source"] = filepath.Base(file)
		metadata["path"] = file

		docs = append(docs, documentSchema.Document{PageContent: text, Metadata: metadata})
	}

	return docs
}
