package documentLoaders

import (
	"context"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	tokenURL = "https://accounts.google.com/o/oauth2/token"
)

type GoogleDriveLoader struct {
	ServiceAccountKey string
	CredentialsPath   string
	TokenPath         string
	FolderID          string
	DocumentIDs       []string
	FileIDs           []string
	Recursive         bool
}

func NewGoogleDriveLoader() *GoogleDriveLoader {
	home := os.Getenv("HOME")
	return &GoogleDriveLoader{
		ServiceAccountKey: filepath.Join(home, ".credentials", "keys.json"),
		CredentialsPath:   filepath.Join(home, ".credentials", "credentials.json"),
		TokenPath:         filepath.Join(home, ".credentials", "token.json"),
	}
}

func (g *GoogleDriveLoader) LoadCredentials() (*http.Client, error) {
	b, err := os.ReadFile(g.ServiceAccountKey)
	if err != nil {
		return nil, err
	}
	// You can add multiple scopes by appending to the second parameter
	config, err := google.JWTConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return nil, err
	}
	client := config.Client(oauth2.NoContext)
	return client, nil
}

func (g *GoogleDriveLoader) LoadSheetFromID(id string) ([]documentSchema.Document, error) {
	var documents []documentSchema.Document
	client, err := g.LoadCredentials()
	if err != nil {
		return []documentSchema.Document{}, err
	}
	service, err := sheets.New(client)
	if err != nil {
		return nil, err
	}
	response, err := service.Spreadsheets.Get(id).Do()
	if err != nil {
		return nil, err
	}
	for _, sheet := range response.Sheets {
		documents = append(documents, documentSchema.Document{PageContent: sheet.Properties.Title, Metadata: map[string]interface{}{"id": response.Properties.Title}})
	}
	return documents, nil
}

func (g *GoogleDriveLoader) LoadDocumentFromID(id string) (documentSchema.Document, error) {
	client, err := g.LoadCredentials()
	if err != nil {
		return documentSchema.Document{}, err
	}
	service, err := drive.New(client)
	if err != nil {
		return documentSchema.Document{}, err
	}
	file, err := service.Files.Get(id).SupportsAllDrives(true).Do()
	if err != nil {
		return documentSchema.Document{}, err
	}
	return documentSchema.Document{PageContent: file.Name, Metadata: map[string]interface{}{"id": file.Id}}, nil
}

func (g *GoogleDriveLoader) LoadDocumentsFromFolder(folderID string) ([]documentSchema.Document, error) {
	client, err := g.LoadCredentials()
	if err != nil {
		return []documentSchema.Document{}, err
	}
	service, err := drive.New(client)
	if err != nil {
		return nil, err
	}
	request := service.Files.List().Q("'" + folderID + "' in parents").IncludeItemsFromAllDrives(true).SupportsAllDrives(true)
	var documents []documentSchema.Document
	err = request.Pages(context.Background(), func(page *drive.FileList) error {
		for _, file := range page.Files {
			if file.MimeType == "application/vnd.google-apps.document" {
				document, err := g.LoadDocumentFromID(file.Id)
				if err != nil {
					return err
				}
				documents = append(documents, document)
			} else if file.MimeType == "application/vnd.google-apps.spreadsheet" {
				document, err := g.LoadSheetFromID(file.Id)
				if err != nil {
					return err
				}
				documents = append(documents, document...)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (g *GoogleDriveLoader) LoadDocumentsFromIDs() ([]documentSchema.Document, error) {
	var documents []documentSchema.Document
	for _, documentID := range g.DocumentIDs {
		document, err := g.LoadDocumentFromID(documentID)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	return documents, nil
}

func (g *GoogleDriveLoader) LoadFileFromID(id string) ([]documentSchema.Document, error) {
	client, err := g.LoadCredentials()
	if err != nil {
		return []documentSchema.Document{}, err
	}
	service, err := drive.New(client)
	if err != nil {
		return nil, err
	}
	request := service.Files.Get(id).SupportsAllDrives(true)
	resp, err := request.Download()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	content := string(body)
	return []documentSchema.Document{documentSchema.Document{PageContent: content, Metadata: map[string]interface{}{"id": id}}}, nil
}

func (g *GoogleDriveLoader) LoadFileFromIDs() ([]documentSchema.Document, error) {
	var documents []documentSchema.Document
	for _, fileID := range g.FileIDs {
		document, err := g.LoadFileFromID(fileID)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document...)
	}
	return documents, nil
}

func (g *GoogleDriveLoader) Load() ([]documentSchema.Document, error) {
	if g.FolderID != "" {
		return g.LoadDocumentsFromFolder(g.FolderID)
	} else if g.DocumentIDs != nil {
		return g.LoadDocumentsFromIDs()
	} else {
		return g.LoadFileFromIDs()
	}
}
