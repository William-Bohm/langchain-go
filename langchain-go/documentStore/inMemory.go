package documentStore

import (
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
)

type InMemoryDocstore struct {
	dict map[string]*documentSchema.Document
}

func NewInMemoryDocstore(dict map[string]*documentSchema.Document) *InMemoryDocstore {
	return &InMemoryDocstore{dict: dict}
}

func (ds *InMemoryDocstore) Add(texts map[string]*documentSchema.Document) {
	overlapping := make(map[string]*documentSchema.Document)
	for k, v := range texts {
		if _, ok := ds.dict[k]; ok {
			overlapping[k] = v
		}
	}
	if len(overlapping) > 0 {
		fmt.Printf("Tried to add ids that already exist: %v\n", overlapping)
		return
	}
	for k, v := range texts {
		ds.dict[k] = v
	}
}

func (ds *InMemoryDocstore) Search(search string) (string, *documentSchema.Document) {
	if val, ok := ds.dict[search]; ok {
		return "", val
	}
	return fmt.Sprintf("ID %s not found.", search), nil
}
