package main

import (
	"fmt"
	"github.com/dfsoccer/linguist/unstructured/partition/html"
)

type UnstructuredHTMLLoader struct {
	UnstructuredFileLoader
}

func (l *UnstructuredHTMLLoader) GetElements() []*html.Element {
	elements, err := html.PartitionHTML(l.FilePath, l.UnstructuredKwargs)
	if err != nil {
		panic(err)
	}
	return elements
}

type UnstructuredFileLoader interface {
	GetElements() []*html.Element
}

func main() {
	// Example usage
	loader := &UnstructuredHTMLLoader{
		UnstructuredFileLoader: UnstructuredFileLoaderImpl{FilePath: "example.html"},
	}
	elements := loader.GetElements()
	for _, element := range elements {
		fmt.Println("Element:", element)
	}
}

type UnstructuredFileLoaderImpl struct {
	FilePath           string
	UnstructuredKwargs map[string]interface{}
}

func main() {
	// Example usage
	loader := &UnstructuredHTMLLoader{
		UnstructuredFileLoader: UnstructuredFileLoaderImpl{FilePath: "example.html"},
	}
	elements := loader.GetElements()
	for _, element := range elements {
		fmt.Println("Element:", element)
	}
}
