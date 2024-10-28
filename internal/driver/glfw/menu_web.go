//go:build wasm || test_web_driver

package glfw

import "djpken/go-fyne"

func addMissingQuitForMainMenu(menus *fyne.MainMenu, w *window) *fyne.MainMenu {
	// no-op for a web browser
	return menus
}
