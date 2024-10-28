//go:build !hints

package theme

import (
	"image/color"

	"github.com/djpken/go-fyne"
)

var (
	fallbackColor = color.Transparent
	fallbackIcon  = &fyne.StaticResource{}
)
