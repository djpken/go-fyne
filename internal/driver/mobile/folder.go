package mobile

import (
	"errors"

	"github.com/djpken/go-fyne"
)

type lister struct {
	fyne.URI
}

func (l *lister) List() ([]fyne.URI, error) {
	return listURI(l)
}

func listerForURI(uri fyne.URI) (fyne.ListableURI, error) {
	if !canListURI(uri) {
		return nil, errors.New("specified URI is not listable")
	}

	return &lister{uri}, nil
}
