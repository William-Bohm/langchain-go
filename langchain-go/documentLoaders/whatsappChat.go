package documentLoaders

import (
	"bufio"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"os"
	"regexp"
	"strings"
)

func ConcatenateRows(date string, sender string, text string) string {
	return fmt.Sprintf("%s on %s: %s\n\n", sender, date, text)
}

type WhatsAppChatLoader struct {
	FilePath string
}

func (w *WhatsAppChatLoader) Load() []documentSchema.Document {
	file, err := os.Open(w.FilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	textContent := ""
	messageLineRegex := regexp.MustCompile("(\\d{1,2}/\\d{1,2}/\\d{2,4}, \\d{1,2}:\\d{1,2}[ _]?(?:AM|PM)?) - (.*?): (.*)")
	for _, line := range lines {
		result := messageLineRegex.FindStringSubmatch(strings.TrimSpace(line))
		if len(result) > 0 {
			date, sender, text := result[1], result[2], result[3]
			textContent += ConcatenateRows(date, sender, text)
		}
	}

	metadata := map[string]interface{}{"source": w.FilePath}

	return []documentSchema.Document{{PageContent: textContent, Metadata: metadata}}
}

func NewWhatsAppChatLoader(path string) *WhatsAppChatLoader {
	return &WhatsAppChatLoader{FilePath: path}
}
