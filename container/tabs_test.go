package container

import (
	"testing"

	internalTest "github.com/djpken/go-fyne/internal/test"
	"github.com/stretchr/testify/assert"

	"github.com/djpken/go-fyne"
	"github.com/djpken/go-fyne/canvas"
	"github.com/djpken/go-fyne/internal/cache"
	"github.com/djpken/go-fyne/test"
	"github.com/djpken/go-fyne/theme"
	"github.com/djpken/go-fyne/widget"
)

func TestTabButton_Icon_Change(t *testing.T) {
	b := &tabButton{icon: theme.CancelIcon()}
	r := cache.Renderer(b)
	icon := r.Objects()[3].(*canvas.Image)
	oldResource := icon.Resource

	b.icon = theme.ConfirmIcon()
	b.Refresh()
	assert.NotEqual(t, oldResource, icon.Resource)
}

func TestTab_ThemeChange(t *testing.T) {
	a := test.NewTempApp(t)
	a.Settings().SetTheme(internalTest.LightTheme(theme.DefaultTheme()))

	tabs := NewAppTabs(
		NewTabItem("a", widget.NewLabel("a")),
		NewTabItem("b", widget.NewLabel("b")))
	w := test.NewTempWindow(t, tabs)
	w.Resize(fyne.NewSize(180, 120))

	initial := w.Canvas().Capture()

	a.Settings().SetTheme(internalTest.DarkTheme(theme.DefaultTheme()))
	tabs.SelectIndex(1)
	second := w.Canvas().Capture()
	assert.NotEqual(t, initial, second)

	a.Settings().SetTheme(internalTest.LightTheme(theme.DefaultTheme()))
	tabs.SelectIndex(0)
	assert.Equal(t, initial, w.Canvas().Capture())
}
