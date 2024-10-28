package common

import (
	"djpken/go-fyne"
	"djpken/go-fyne/internal/cache"
)

// CanvasForObject returns the canvas for the specified object.
func CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	return cache.GetCanvasForObject(obj)
}
