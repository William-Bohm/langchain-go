package loading

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	defaultRef = os.Getenv("LANGCHAIN_HUB_DEFAULT_REF")
	urlBase    = os.Getenv("LANGCHAIN_HUB_URL_BASE")
	hubPathRe  = regexp.MustCompile(`lc(?P<ref>@[^:]+)?://(?P<path>.*)`)
)

func TryLoadFromHub(
	path string,
	loader func(string, map[string]interface{}) (interface{}, error),
	validPrefix string,
	validSuffixes []string,
	kwargs map[string]interface{},
) (interface{}, error) {
	if _, err := url.ParseRequestURI(path); err != nil || !hubPathRe.MatchString(path) {
		return nil, nil
	}

	matches := hubPathRe.FindStringSubmatch(path)
	ref := matches[1]
	if ref == "" {
		ref = defaultRef
	} else {
		ref = ref[1:]
	}
	remotePathStr := matches[2]
	remotePath := filepath.Clean(remotePathStr)
	if strings.Split(remotePath, "/")[0] != validPrefix {
		return nil, nil
	}
	if !contains(validSuffixes, filepath.Ext(remotePath)) {
		return nil, fmt.Errorf("Unsupported file type.")
	}

	fullURL := urlBase + "/" + ref + "/" + remotePath

	resp, err := http.Get(fullURL)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("Could not find file at %s", fullURL)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	tmpDirName := os.TempDir() + "/" + uuid.New().String()
	os.Mkdir(tmpDirName, 0755)
	defer os.RemoveAll(tmpDirName)

	file := tmpDirName + "/" + filepath.Base(remotePath)
	os.WriteFile(file, body, 0644)

	return loader(file, kwargs), nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
