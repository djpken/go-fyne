package mobile

import (
	"testing"

	_ "djpken/go-fyne/test"
)

func TestDevice_HideVirtualKeyboard_BeforeRun(t *testing.T) {
	hideVirtualKeyboard() // should not crash!
}
