//go:build ios

package mobile

import (
	fyneDriver "github.com/djpken/go-fyne/driver"
)

// Assert we are satisfying the driver.NativeWindow interface
var _ fyneDriver.NativeWindow = (*window)(nil)

func (w *window) RunNative(fn func(context any)) {
	fn(&fyneDriver.UnknownContext{})
}
