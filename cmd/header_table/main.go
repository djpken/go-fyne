package main

import (
	"djpken/go-fyne"
	"djpken/go-fyne/app"
	"djpken/go-fyne/container"
	"djpken/go-fyne/widget"
	"fmt"
)

func main() {
	app := app.New()
	window := app.NewWindow("test")
	tab := makeTableTab(window)
	window.SetContent(tab)
	window.Show()
	window.Resize(fyne.Size{Height: 400, Width: 600})
	app.Run()

}

func makeTableTab(_ fyne.Window) fyne.CanvasObject {
	t := widget.NewStaticTable(
		func() (int, int) { return 40, 20 },
		func() fyne.CanvasObject {
			return container.NewStack()
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := widget.NewLabel("")
			label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
			cell.(*fyne.Container).Add(label)
		})
	t.SetColumnWidth(0, 102)
	return t
}
