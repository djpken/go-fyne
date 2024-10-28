package dialog

import (
	"testing"

	"djpken/go-fyne"
	"djpken/go-fyne/test"
)

func Test_colorWheel_Layout(t *testing.T) {
	test.NewTempApp(t)

	wheel := newColorWheel(nil)
	wheel.SetHSLA(180, 100, 50, 255)
	window := test.NewTempWindow(t, wheel)
	window.Resize(wheel.MinSize().Max(fyne.NewSize(100, 100)))

	test.AssertRendersToImage(t, "color/wheel_layout.png", window.Canvas())
	test.AssertRendersToMarkup(t, "color/wheel_layout.xml", window.Canvas())
}
