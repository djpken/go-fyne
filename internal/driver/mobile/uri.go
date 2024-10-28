//go:build !android

package mobile

import (
	"github.com/djpken/go-fyne"
	"github.com/djpken/go-fyne/storage"
)

func nativeURI(path string) fyne.URI {
	uri, err := storage.ParseURI(path)
	if err != nil {
		fyne.LogError("Error on parsing uri", err)
	}
	return uri
}
