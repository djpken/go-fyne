//go:build ci

package app

import (
	"djpken/go-fyne"
	"djpken/go-fyne/internal/painter/software"
	"djpken/go-fyne/test"
)

// NewWithID returns a new app instance using the test (headless) driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	return newAppWithDriver(test.NewDriverWithPainter(software.NewPainter()), id)
}
