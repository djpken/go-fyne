package widget

import (
	"testing"

	"djpken/go-fyne/canvas"
	"djpken/go-fyne/test"
	"djpken/go-fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestNewIcon(t *testing.T) {
	icon := NewIcon(theme.ConfirmIcon())
	render := test.TempWidgetRenderer(t, icon)

	assert.Equal(t, 1, len(render.Objects()))
	obj := render.Objects()[0]
	img, ok := obj.(*canvas.Image)
	if !ok {
		t.Fail()
	}
	assert.Equal(t, theme.ConfirmIcon(), img.Resource)
}

func TestIcon_Nil(t *testing.T) {
	icon := NewIcon(nil)
	render := test.TempWidgetRenderer(t, icon)

	assert.Equal(t, 1, len(render.Objects()))
	assert.Nil(t, render.Objects()[0].(*canvas.Image).Resource)
}

func TestIcon_MinSize(t *testing.T) {
	icon := NewIcon(theme.CancelIcon())
	min := icon.MinSize()

	assert.Equal(t, theme.IconInlineSize(), min.Width)
	assert.Equal(t, theme.IconInlineSize(), min.Height)
}

func TestIconRenderer_ApplyTheme(t *testing.T) {
	icon := NewIcon(theme.CancelIcon())
	render := test.TempWidgetRenderer(t, icon).(*iconRenderer)
	visible := render.Objects()[0].Visible()

	render.Refresh()
	assert.Equal(t, visible, render.Objects()[0].Visible())
}
