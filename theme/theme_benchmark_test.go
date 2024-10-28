package theme

import (
	"testing"

	"djpken/go-fyne"
)

func BenchmarkTheme_current(b *testing.B) {
	fyne.CurrentApp().Settings().SetTheme(LightTheme())

	for n := 0; n < b.N; n++ {
		Current()
	}
}
