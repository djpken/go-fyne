// Package main loads a very basic Hello World graphical application.
package main

import (
	"djpken/go-fyne/app"
	"djpken/go-fyne/container"
	"djpken/go-fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome 😀")
		}),
	))

	w.ShowAndRun()
}
