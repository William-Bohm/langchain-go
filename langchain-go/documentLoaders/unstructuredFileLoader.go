package documentLoaders

import (
	"errors"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"strings"
)

type UnstructuredBaseLoader struct {
	Mode               string
	UnstructuredKwargs map[string]interface{}
}

type UnstructuredFileLoader struct {
	UnstructuredBaseLoader
	FilePath string
}

type UnstructuredFileIOLoader struct {
	UnstructuredBaseLoader
	File interface{}
}

func NewUnstructuredFileLoader(filePath string, mode string, unstructuredKwargs map[string]interface{}) UnstructuredFileLoader {
	return UnstructuredFileLoader{
		UnstructuredBaseLoader: UnstructuredBaseLoader{
			Mode:               mode,
			UnstructuredKwargs: unstructuredKwargs,
		},
		FilePath: filePath,
	}
}

func SatisfiesMinUnstructuredVersion(minVersion string) bool {
	unstructuredVersion := "0.6.0"
	minVersionElements := strings.Split(minVersion, ".")
	unstructuredVersionElements := strings.Split(unstructuredVersion, ".")
	for i := 0; i < len(minVersionElements); i++ {
		if minVersionElements[i] < unstructuredVersionElements[i] {
			return false
		}
	}
	return true
}

func (u *UnstructuredBaseLoader) Load() ([]documentSchema.Document, error) {
	return nil, errors.New("not implemented")
}

func (u *UnstructuredFileLoader) GetElements() ([]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (u *UnstructuredFileLoader) GetMetadata() (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (u *UnstructuredFileIOLoader) GetElements() ([]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (u *UnstructuredFileIOLoader) GetMetadata() (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}
