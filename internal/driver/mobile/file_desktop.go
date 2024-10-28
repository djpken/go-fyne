//go:build !ios && !android

package mobile

import (
	"io"

	"djpken/go-fyne"
	intRepo "djpken/go-fyne/internal/repository"
	"djpken/go-fyne/storage/repository"
)

func deleteURI(u fyne.URI) error {
	// no-op as we use the internal FileRepository
	return nil
}

func existsURI(fyne.URI) (bool, error) {
	// no-op as we use the internal FileRepository
	return false, nil
}

func nativeFileOpen(*fileOpen) (io.ReadCloser, error) {
	// no-op as we use the internal FileRepository
	return nil, nil
}

func nativeFileSave(*fileSave) (io.WriteCloser, error) {
	// no-op as we use the internal FileRepository
	return nil, nil
}

func registerRepository(d *driver) {
	repository.Register("file", intRepo.NewFileRepository())
}
