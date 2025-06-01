// TODO main.go sounds like mango if you say it fast LOL

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Board struct {
	Title   string   `json:"title"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Title string `json:"title"`
	Cards []Card `json:"cards"`
}

type Card struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Shortcut    string `json:"shortcut"`
}

var boardData Board

var columns []*tview.List
var currentFocus int = 0

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

func loadData() {
	filePath := "./tasks.json" // TODO set this up as a cli input to potentially load different boards from different JSONs

	var file *os.File
	var err error
	if !fileExists(filePath) {
		file, err = os.Create(filePath)
	} else {
		file, err = os.Open("tasks.json")
	}

	if err != nil {
		panic(err)
	}

	defer file.Close()
	json.NewDecoder(file).Decode(&boardData)
}

func setBoard(app *tview.Application, board *tview.Flex) {
	for i, col := range boardData.Columns {
		list := tview.NewList()
		list.SetBorder(true)
		list.SetTitle(col.Title)

		for _, card := range col.Cards {
			shortcutRune := []rune(card.Shortcut)

			var r rune
			if len(shortcutRune) > 0 {
				r = shortcutRune[0]
			} else {
				r = 0
			}

			list.AddItem(card.Title, card.Description, r, func() {
				modal := tview.NewModal().SetBackgroundColor(tcell.Color158).SetText(fmt.Sprintf("Title: %s\n\nDescription: %s", card.Title, card.Description)).AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetRoot(board, true).SetFocus(list)
				})
				modal.SetBorder(false)
				app.SetRoot(modal, true).SetFocus(modal)
			})
		}

		board.AddItem(list, 0, 1, i == 0)

		columns = append(columns, list)
	}
}

func main() {
	loadData()

	app := tview.NewApplication()
	board := tview.NewFlex().SetDirection(tview.FlexColumn)

	setBoard(app, board)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			addCard(columns[currentFocus])
		case 'd':
			deleteCard(columns[currentFocus])
		case 'h':
			if currentFocus > 0 {
				currentFocus--
				app.SetFocus(columns[currentFocus])
			}
		case 'j':
			list := columns[currentFocus]
			index := list.GetCurrentItem()
			if index < list.GetItemCount()-1 {
				list.SetCurrentItem(index + 1)
			} else {
				list.SetCurrentItem(0)
			}
		case 'k':
			list := columns[currentFocus]
			index := list.GetCurrentItem()
			if index > 0 {
				list.SetCurrentItem(index - 1)
			} else {
				list.SetCurrentItem(list.GetItemCount() - 1)
			}
		case 'l':
			if currentFocus < len(columns)-1 {
				currentFocus++
				app.SetFocus(columns[currentFocus])
			}
		case 'p':
			promoteCard(app, board)
		case 'q':
			app.Stop()
		case 'r':
			regressCard(app, board)
		}
		return event
	})

	if err := app.SetRoot(board, true).SetFocus(columns[0]).Run(); err != nil {
		panic(err)
	}
}

func deleteCard(column *tview.List) {
	currentItemIndex := column.GetCurrentItem()
	column.RemoveItem(currentItemIndex)
}

func addCard(column *tview.List) {
	column.AddItem("HARDCODED ADD", "DETAILED DESCRIPTION", '1', nil)
}

func promoteCard(app *tview.Application, board *tview.Flex) {
	if currentFocus >= len(columns)-1 {
		return
	}

	adjacentList := columns[currentFocus+1]

	list := columns[currentFocus]
	if list.GetItemCount() == 0 {
		return
	}

	index := list.GetCurrentItem()

	cards := boardData.Columns[currentFocus].Cards
	card := cards[index]
	boardData.Columns[currentFocus].Cards = slices.Delete(cards, index, index+1)

	list.RemoveItem(index)

	shortcutRune := []rune(card.Shortcut)

	var r rune
	if len(shortcutRune) > 0 {
		r = shortcutRune[0]
	} else {
		r = 0
	}

	adjacentList.AddItem(card.Title, card.Description, r, func() {
		modal := tview.NewModal().SetBackgroundColor(tcell.Color158).SetText(fmt.Sprintf("Title: %s\n\nDescription: %s", card.Title, card.Description)).AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(board, true).SetFocus(list)
		})
		modal.SetBorder(false)
		app.SetRoot(modal, true).SetFocus(modal)
	})

	boardData.Columns[currentFocus+1].Cards = append(boardData.Columns[currentFocus+1].Cards, card)
}

func regressCard(app *tview.Application, board *tview.Flex) {
	if currentFocus == 0 {
		return
	}

	adjacentList := columns[currentFocus-1]

	list := columns[currentFocus]
	if list.GetItemCount() == 0 {
		return
	}

	index := list.GetCurrentItem()

	cards := boardData.Columns[currentFocus].Cards
	card := cards[index]
	boardData.Columns[currentFocus].Cards = slices.Delete(cards, index, index+1)

	list.RemoveItem(index)

	shortcutRune := []rune(card.Shortcut)

	var r rune
	if len(shortcutRune) > 0 {
		r = shortcutRune[0]
	} else {
		r = 0
	}

	adjacentList.AddItem(card.Title, card.Description, r, func() {
		modal := tview.NewModal().SetBackgroundColor(tcell.Color158).SetText(fmt.Sprintf("Title: %s\n\nDescription: %s", card.Title, card.Description)).AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(board, true).SetFocus(list)
		})
		modal.SetBorder(false)
		app.SetRoot(modal, true).SetFocus(modal)
	})

	boardData.Columns[currentFocus-1].Cards = append(boardData.Columns[currentFocus-1].Cards, card)
}
