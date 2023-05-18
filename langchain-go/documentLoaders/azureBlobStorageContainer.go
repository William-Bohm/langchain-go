package documentLoaders

import (
	"context"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"net/url"
	"os"
	"strings"
)

type AzureBlobStorageContainerLoader struct {
	BaseLoaderImpl
	ConnStr   string
	Container string
	Prefix    string
}

func NewAzureBlobStorageContainerLoader(connStr string, container string, prefix string) *AzureBlobStorageContainerLoader {
	return &AzureBlobStorageContainerLoader{
		ConnStr:   connStr,
		Container: container,
		Prefix:    prefix,
	}
}

func (a *AzureBlobStorageContainerLoader) Load() ([]documentSchema.Document, error) {
	accountName, _ := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey, _ := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	credential, _ := azblob.NewSharedKeyCredential(accountName, accountKey)
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(
		"https://" + accountName + ".blob.core.windows.net/" + a.Container + "/")
	containerURL := azblob.NewContainerURL(*URL, p)
	ctx := context.Background()
	listBlobs, _ := containerURL.ListBlobsFlatSegment(ctx, azblob.Marker{}, azblob.ListBlobsSegmentOptions{})

	var docs []documentSchema.Document
	for _, blobInfo := range listBlobs.Segment.BlobItems {
		if strings.HasPrefix(blobInfo.Name, a.Prefix) {
			loader := NewAzureBlobStorageFileLoader(a.ConnStr, a.Container, blobInfo.Name)
			loadedDocs, _ := loader.Load()
			docs = append(docs, loadedDocs...)
		}
	}
	return docs, nil
}
