package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"log"
	"math/rand/v2"
	"time"
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
	t := widget.NewStaticTable(
		func() (int, int) { return len(bindings), 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			var str any
			switch i.Col {
			case 0:
				str, _ = bindings[i.Row].GetValue("UserID")
			case 1:
				str, _ = bindings[i.Row].GetValue("ID")
			case 2:
				str, _ = bindings[i.Row].GetValue("Title")
			case 3:
				str, _ = bindings[i.Row].GetValue("Completed")
			}
			o.(*widget.Label).SetText(str.(string))
		})
	t.SetColumnWidth(0, 102)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				i := rand.IntN(1000)
				bindings = append(bindings, binding.BindStruct(&Todo{
					UserID:    "1",
					ID:        "4",
					Title:     "et porro tempora",
					Completed: "true",
				}))
				bindings[1] = binding.BindStruct(&Todo{
					UserID:    "1",
					ID:        fmt.Sprintf("%d", i),
					Title:     "et porro tempora",
					Completed: "true",
				})
				fyne.Do(func() {
					t.Refresh()
				})
			}
		}
	}()

	return t
}
