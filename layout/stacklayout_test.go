package layout_test

import (
	"image/color"
	"testing"

	"djpken/go-fyne"
	"djpken/go-fyne/canvas"
	"djpken/go-fyne/container"
	"djpken/go-fyne/layout"

	"github.com/stretchr/testify/assert"
)

func TestStackLayout(t *testing.T) {
	size := fyne.NewSize(100, 100)

	obj := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj},
	}
	container.Resize(size)

	layout.NewStackLayout().Layout(container.Objects, size)

	assert.Equal(t, obj.Size(), size)
}

func TestStackLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := container.NewWithoutLayout(text)
	layoutMin := layout.NewStackLayout().MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestContainerStackLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.NRGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := container.NewWithoutLayout(text)
	container.Layout = layout.NewStackLayout()
	layoutMin := container.MinSize()

	assert.Equal(t, minSize, layoutMin)
}
