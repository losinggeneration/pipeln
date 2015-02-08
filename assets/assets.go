package assets

import (
	"errors"
	"os"
	"path/filepath"
)

type Options struct {
	BaseDir        string
	JavascriptsDir string
	StylesheetsDir string
	ImagesDir      string
}

type Assets struct {
	Options Options
}

var (
	ErrNotFound = errors.New("asset not found")
)

func (a Assets) Lookup(name string) (*os.File, error) {
	path := filepath.Join(a.Options.BaseDir, name)
	file, err := os.Open(path)

	if err != nil {
		return nil, ErrNotFound
	}

	return file, nil
}

func (a Assets) Process(name string) error {
	return nil
}
