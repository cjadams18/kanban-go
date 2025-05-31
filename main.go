package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO dynamically add tasks
// TODO store tasks in SQLite db? Or JSON (maybe JSON so shareable extract)
// TODO use vim controls to navigate, but also make something for arrow keys as well
func main() {
	app := tview.NewApplication()

	backlogColumn := tview.NewList()
	backlogColumn.SetBorder(true)
	backlogColumn.SetTitle(" Backlog ")

	backlogColumn.AddItem("Task 1", "This is a description", '1', nil)
	backlogColumn.AddItem("Task 2", "THIS IS A TEST", '2', nil)

	inProgressColumn := tview.NewList()
	inProgressColumn.SetBorder(true)
	inProgressColumn.SetTitle(" In progress ")

	inProgressColumn.AddItem("Task 2", "Another very detailed task todo", '1', nil)

	doneColumn := tview.NewList()
	doneColumn.SetBorder(true)
	doneColumn.SetTitle(" Ready for Test ")

	columns := []*tview.List{backlogColumn, inProgressColumn, doneColumn}
	currentFocus := 0

	board := tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(backlogColumn, 0, 1, true).AddItem(inProgressColumn, 0, 1, false).AddItem(doneColumn, 0, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'h':
			if currentFocus > 0 {
				currentFocus--
				app.SetFocus(columns[currentFocus])
			}
		case 'j':
		case 'k':
		case 'l':
			if currentFocus < len(columns)-1 {
				currentFocus++
				app.SetFocus(columns[currentFocus])
			}
		case 'q':
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(board, true).SetFocus(backlogColumn).Run(); err != nil {
		panic(err)
	}
}
