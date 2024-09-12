package examples

import (
	"embed"
	"errors"
	"gopkg.in/yaml.v2"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:embed
var examplesFS embed.FS

var ErrExampleNotFound = errors.New("example not found")

type Example struct {
	Name        string
	Title       string
	Description string
}

func List() ([]Example, error) {
	entries, err := examplesFS.ReadDir("embed/dist")
	if err != nil {
		return nil, err
	}

	exampleList := make([]Example, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		rillYamlContents, err := examplesFS.ReadFile(filepath.Join("embed/dist", entry.Name(), "rill.yaml"))
		if err != nil {
			return nil, err
		}

		contents := struct {
			Title       string
			Description string
		}{}

		if err := yaml.Unmarshal(rillYamlContents, &contents); err != nil {
			return nil, err
		}

		exampleList = append(exampleList, Example{
			Name:        entry.Name(),
			Title:       contents.Title,
			Description: contents.Description,
		})
	}

	return exampleList, nil
}

func Get(name string) (fs.FS, error) {
	exampleFS, err := fs.Sub(examplesFS, filepath.Join("embed", "dist", name))
	if err != nil {
		return nil, err
	}

	_, err = fs.Stat(exampleFS, "rill.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrExampleNotFound
		}
		return nil, err
	}

	return exampleFS, nil
}
