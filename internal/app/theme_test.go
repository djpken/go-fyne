package app_test

import (
	"testing"

	"github.com/djpken/go-fyne/internal/app"
	"github.com/djpken/go-fyne/test"
)

func TestApplySettings_BeforeContentSet(t *testing.T) {
	a := test.NewApp()
	w := a.NewWindow("NoContent")
	defer w.Close()

	app.ApplySettings(a.Settings(), a)
}
