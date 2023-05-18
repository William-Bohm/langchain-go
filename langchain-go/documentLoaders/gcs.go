package documentLoaders

import (
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/ioutil"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

type GCSDirectoryLoader struct {
	ProjectName string
	Bucket      string
	Prefix      string
}

func (g *GCSDirectoryLoader) Load() ([]documentSchema.Document, error) {
	ctx := context.Background()

	client, _ := storage.NewClient(ctx)
	bkt := client.Bucket(g.Bucket)
	query := &storage.Query{Prefix: g.Prefix}
	it := bkt.Objects(ctx, query)

	docs := []documentSchema.Document{}

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		loader := GCSFileLoader{
			ProjectName: g.ProjectName,
			Bucket:      g.Bucket,
			Blob:        attrs.Name,
		}
		extraDocs, err := loader.Load()
		if err != nil {
			return []documentSchema.Document{}, err
		}
		docs = append(docs, extraDocs...)
	}
	return docs, nil
}

type GCSFileLoader struct {
	ProjectName string
	Bucket      string
	Blob        string
}

func (g *GCSFileLoader) Load() ([]documentSchema.Document, error) {
	ctx := context.Background()

	client, _ := storage.NewClient(ctx)
	bkt := client.Bucket(g.Bucket)
	blob := bkt.Object(g.Blob)

	r, _ := blob.NewReader(ctx)
	defer r.Close()

	data, _ := ioutil.ReadAll(r)

	tmpDir, _ := ioutil.TempDir("", "example")
	defer os.RemoveAll(tmpDir) // clean up

	tmpFileName := filepath.Join(tmpDir, g.Blob)
	os.WriteFile(tmpFileName, data, 0644)

	loader := UnstructuredFileLoader{
		FilePath: tmpFileName,
	}
	return loader.Load()
}
