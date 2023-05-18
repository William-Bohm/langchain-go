package documentLoaders

import (
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type FILE_LOADER_TYPE interface {
	Load() []documentSchema.Document
}

type DirectoryLoader struct {
	path         string
	glob         string
	loadHidden   bool
	loaderCls    FILE_LOADER_TYPE
	loaderKwargs map[string]interface{}
	silentErrors bool
	recursive    bool
}

func isVisible(p string) bool {
	_, file := filepath.Split(p)
	return !strings.HasPrefix(file, ".")
}

func NewDirectoryLoader(path string, glob string, silentErrors bool, loadHidden bool, loaderCls FILE_LOADER_TYPE, loaderKwargs map[string]interface{}, recursive bool) *DirectoryLoader {
	if loaderKwargs == nil {
		loaderKwargs = make(map[string]interface{})
	}

	return &DirectoryLoader{
		path:         path,
		glob:         glob,
		loadHidden:   loadHidden,
		loaderCls:    loaderCls,
		loaderKwargs: loaderKwargs,
		silentErrors: silentErrors,
		recursive:    recursive,
	}
}

func (d *DirectoryLoader) Load() []documentSchema.Document {
	var docs []documentSchema.Document

	err := filepath.Walk(d.path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isVisible(path) || d.loadHidden {
			subDocs := d.loaderCls.Load()
			docs = append(docs, subDocs...)
		}
		return nil
	})

	if err != nil {
		if !d.silentErrors {
			log.Fatal(err)
		} else {
			log.Println(err)
		}
	}

	return docs
}
