//go:build !hints

package theme

import (
	"image/color"

	"djpken/go-fyne"
)

var (
	fallbackColor = color.Transparent
	fallbackIcon  = &fyne.StaticResource{}
)
