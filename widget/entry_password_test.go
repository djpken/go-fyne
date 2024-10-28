package widget_test

import (
	"testing"

	"github.com/djpken/go-fyne"
	"github.com/djpken/go-fyne/canvas"
	"github.com/djpken/go-fyne/container"
	"github.com/djpken/go-fyne/test"
	"github.com/djpken/go-fyne/widget"
	"github.com/stretchr/testify/assert"
)

func TestNewPasswordEntry(t *testing.T) {
	p := widget.NewPasswordEntry()
	p.Text = "visible"
	r := test.TempWidgetRenderer(t, p)

	cont := r.Objects()[2].(*container.Scroll).Content.(fyne.Widget)
	r = test.TempWidgetRenderer(t, cont)
	rich := r.Objects()[1].(*widget.RichText)
	r = test.TempWidgetRenderer(t, rich)

	assert.Equal(t, "•••••••", r.Objects()[0].(*canvas.Text).Text)
}
