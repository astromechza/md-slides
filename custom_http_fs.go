package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type CustomDirFS struct {
	Directory string
}

var _ http.FileSystem = (*CustomDirFS)(nil)

func (fs CustomDirFS) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.Join(fs.Directory, filepath.FromSlash(path.Clean("/"+name)))
	log.Printf("Attempting to load: %s", fullName)

	f, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}
	d, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if d.IsDir() {
		return nil, fmt.Errorf("cannot serve directory")
	}
	return f, nil
}
