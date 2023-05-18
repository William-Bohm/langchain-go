package documentSchema

type Document struct {
	PageContent string                 `json:"page_content"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type Docstore interface {
	Search(search string) (string, *Document)
}

type AddableMixin interface {
	Add(texts map[string]*Document)
}
