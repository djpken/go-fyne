package mobile

import (
	"github.com/djpken/go-fyne"
	"github.com/djpken/go-fyne/driver/mobile"
	"github.com/djpken/go-fyne/internal/driver/mobile/app"
)

func hideVirtualKeyboard() {
	if d, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if d.app == nil { // not yet running
			return
		}

		d.app.HideVirtualKeyboard()
	}
}

func handleKeyboard(obj fyne.Focusable) {
	isDisabled := false
	if disWid, ok := obj.(fyne.Disableable); ok {
		isDisabled = disWid.Disabled()
	}
	if obj != nil && !isDisabled {
		if keyb, ok := obj.(mobile.Keyboardable); ok {
			showVirtualKeyboard(keyb.Keyboard())
		} else {
			showVirtualKeyboard(mobile.DefaultKeyboard)
		}
	} else {
		hideVirtualKeyboard()
	}
}

func showVirtualKeyboard(keyboard mobile.KeyboardType) {
	if d, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if d.app == nil { // not yet running
			fyne.LogError("Cannot show keyboard before app is running", nil)
			return
		}

		d.app.ShowVirtualKeyboard(app.KeyboardType(keyboard))
	}
}
