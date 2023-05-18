package documentLoaders

import (
	"context"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

type AzureBlobStorageFileLoader struct {
	BaseLoaderImpl
	ConnStr   string
	Container string
	Blob      string
}

func NewAzureBlobStorageFileLoader(connStr string, container string, blob string) *AzureBlobStorageFileLoader {
	return &AzureBlobStorageFileLoader{
		ConnStr:   connStr,
		Container: container,
		Blob:      blob,
	}
}

func (a *AzureBlobStorageFileLoader) Load() ([]documentSchema.Document, error) {
	accountName, _ := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey, _ := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	credential, _ := azblob.NewSharedKeyCredential(accountName, accountKey)
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(
		"https://" + accountName + ".blob.core.windows.net/" + a.Container + "/" + a.Blob)
	blobURL := azblob.NewBlockBlobURL(*URL, p)
	ctx := context.Background()
	blobData, _ := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})

	reader := blobData.Body(azblob.RetryReaderOptions{})
	tempDir, _ := ioutil.TempDir("", "azure-storage-")
	filePath := filepath.Join(tempDir, a.Container, a.Blob)
	_ = os.MkdirAll(filepath.Dir(filePath), 0755)
	file, _ := os.Create(filePath)
	_, _ = io.Copy(file, reader)

	loader := UnstructuredFileLoader{FilePath: filePath}
	docs, _ := loader.Load()
	_ = file.Close()
	_ = os.RemoveAll(tempDir)

	return docs, nil
}
