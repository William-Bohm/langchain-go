package documentLoaders

import (
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"path/filepath"
)

type FileFilter func(string) bool

type GitLoader struct {
	RepoPath   string
	CloneURL   string
	Branch     string
	FileFilter FileFilter
}

func (g *GitLoader) Load() []documentSchema.Document {
	var r *git.Repository
	var err error

	if _, err := os.Stat(g.RepoPath); os.IsNotExist(err) && g.CloneURL == "" {
		panic("Path does not exist")
	} else if g.CloneURL != "" {
		r, err = git.PlainClone(g.RepoPath, false, &git.CloneOptions{
			URL:      g.CloneURL,
			Progress: os.Stdout,
		})
		if err != nil {
			panic(err)
		}
	} else {
		r, err = git.PlainOpen(g.RepoPath)
		if err != nil {
			panic(err)
		}
	}

	ref, err := r.Head()
	if err != nil {
		panic(err)
	}

	branchRef := plumbing.NewHashReference(plumbing.ReferenceName(g.Branch), ref.Hash())
	err = r.Storer.SetReference(branchRef)
	if err != nil {
		panic(err)
	}

	iter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}

	docs := []documentSchema.Document{}

	err = iter.ForEach(func(c *object.Commit) error {
		files, err := c.Files()
		if err != nil {
			panic(err)
		}

		files.ForEach(func(f *object.File) error {
			filePath := filepath.Join(g.RepoPath, f.Name)

			if g.FileFilter != nil && !g.FileFilter(filePath) {
				return nil
			}

			data, err := f.Contents()
			if err != nil {
				panic(err)
			}

			fileType := filepath.Ext(f.Name)
			relFilePath, _ := filepath.Rel(g.RepoPath, filePath)

			metadata := map[string]interface{}{
				"file_path": relFilePath,
				"file_name": f.Name,
				"file_type": fileType,
			}

			doc := documentSchema.Document{
				PageContent: data,
				Metadata:    metadata,
			}

			docs = append(docs, doc)
			return nil
		})
		return nil
	})
	if err != nil {
		panic(err)
	}

	return docs
}
