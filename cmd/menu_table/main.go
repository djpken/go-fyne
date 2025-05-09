package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"log"
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

type Todo struct {
	UserID    string `json:"userId,omitempty"`
	ID        string `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Completed string `json:"completed,omitempty"`
}

func makeTableTab(_ fyne.Window) fyne.CanvasObject {
	var data []Todo

	stringData := `[{"userId":"1","id":"1","title":"delectus aut autem","completed":"false"},{"userId":"1","id":"2","title":"quis ut nam facilis et officia qui","completed":"false"},{"userId":"1","id":"3","title":"fugiat veniam minus","completed":"false"},{"userId":"1","id":"4","title":"et porro tempora","completed":"true"}]`
	err := json.Unmarshal([]byte(stringData), &data)
	if err != nil {
		log.Fatal(err)
	}

	var bindings []binding.Struct

	for _, todo := range data {
		bindings = append(bindings, binding.BindStruct(&todo))
	}
	t := widget.NewMenuTable(
		func() (int, int) { return len(bindings), 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			var str any
			switch i.Col {
			case 1:
				str, _ = bindings[i.Row].GetValue("ID")
			case 2:
				str, _ = bindings[i.Row].GetValue("Title")
			case 3:
				str, _ = bindings[i.Row].GetValue("Completed")
			}
			label, ok := o.(*widget.Label)
			if !ok {
				return
			}
			label.SetText(str.(string))
		})
	t.UpdateMenuButton = func(i widget.TableCellID, o fyne.CanvasObject) {
		button := o.(*widget.MenuButton)
		if i.Row == 5 {
			menu := &fyne.Menu{
				Items: []*fyne.MenuItem{
					fyne.NewMenuItem("5", nil),
				},
			}
			button.SetMenu(menu)
			return
		}
		menu := &fyne.Menu{
			Items: []*fyne.MenuItem{
				fyne.NewMenuItem("edit", nil),
			},
		}
		button.SetMenu(menu)
	}

	return t
}
