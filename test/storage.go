package test

import (
	"os"

	"github.com/djpken/go-fyne"
	"github.com/djpken/go-fyne/internal"
	"github.com/djpken/go-fyne/storage"
)

type testStorage struct {
	*internal.Docs
}

func (s *testStorage) RootURI() fyne.URI {
	return storage.NewFileURI(os.TempDir())
}

func (s *testStorage) docRootURI() (fyne.URI, error) {
	return storage.Child(s.RootURI(), "Documents")
}
