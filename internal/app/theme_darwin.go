//go:build !ios && !wasm && !test_web_driver && !mobile

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#include <AppKit/AppKit.h>

bool isDarkMode();
*/
import "C"
import (
	"github.com/djpken/go-fyne"
	"github.com/djpken/go-fyne/theme"
)

// DefaultVariant returns the systems default fyne.ThemeVariant.
// Normally, you should not need this. It is extracted out of the root app package to give the
// settings app access to it.
func DefaultVariant() fyne.ThemeVariant {
	if C.isDarkMode() {
		return theme.VariantDark
	}
	return theme.VariantLight
}
