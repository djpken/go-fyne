package app_test

import (
	"testing"

	"djpken/go-fyne/internal/app"
	"djpken/go-fyne/test"
)

func TestApplySettings_BeforeContentSet(t *testing.T) {
	a := test.NewApp()
	w := a.NewWindow("NoContent")
	defer w.Close()

	app.ApplySettings(a.Settings(), a)
}
