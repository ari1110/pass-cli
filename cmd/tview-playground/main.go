package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	currentPage := 0
	maxPages := 3

	// Header
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	header.SetBackgroundColor(tcell.ColorPurple)

	// Footer
	footer := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("h/left: Previous | l/right: Next | q: Quit")
	footer.SetBackgroundColor(tcell.ColorDarkGray)

	// Content area - will be swapped based on page
	contentArea := tview.NewFlex()

	// Main layout with flex
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, false).
		AddItem(contentArea, 0, 1, true).
		AddItem(footer, 1, 0, false)

	// Declare updatePage
	var updatePage func()

	// Function to update display
	updatePage = func() {
		header.SetText(fmt.Sprintf("[white]Page %d/%d - Resize terminal to test layout", currentPage, maxPages))

		contentArea.Clear()

		switch currentPage {
		case 0:
			// Page 0: Simple content
			simpleBox := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("[yellow]Page 0:[white]\n\nSimple content box")
			simpleBox.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGray)
			contentArea.AddItem(simpleBox, 0, 1, false)

		case 1:
			// Page 1: Two horizontal panels (1:2 ratio)
			leftPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Left Panel\n\nFlex: 1")
			leftPanel.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGray)

			rightPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Right Panel\n\nFlex: 2")
			rightPanel.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGray)

			contentArea.SetDirection(tview.FlexColumn).
				AddItem(leftPanel, 0, 1, false).
				AddItem(rightPanel, 0, 2, false)

		case 2:
			// Page 2: Three horizontal panels (1:2:1 ratio)
			leftPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Left\n\nFlex: 1")
			leftPanel.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGray)

			centerPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Center\n\nFlex: 2")
			centerPanel.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGray)

			rightPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Right\n\nFlex: 1")
			rightPanel.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGray)

			contentArea.SetDirection(tview.FlexColumn).
				AddItem(leftPanel, 0, 1, false).
				AddItem(centerPanel, 0, 2, false).
				AddItem(rightPanel, 0, 1, false)

		case 3:
			// Page 3: Vertical panels (1:3:1 ratio)
			topPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Top\nFlex: 1")
			topPanel.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGray)

			middlePanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Middle\nFlex: 3")
			middlePanel.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGray)

			bottomPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Bottom\nFlex: 1")
			bottomPanel.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGray)

			contentArea.SetDirection(tview.FlexRow).
				AddItem(topPanel, 0, 1, false).
				AddItem(middlePanel, 0, 3, false).
				AddItem(bottomPanel, 0, 1, false)
		}
	}

	// Initial display
	updatePage()

	// Global input handler
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				app.Stop()
				return nil
			case 'h':
				if currentPage > 0 {
					currentPage--
					updatePage()
				}
				return nil
			case 'l':
				if currentPage < maxPages {
					currentPage++
					updatePage()
				}
				return nil
			}
		case tcell.KeyLeft:
			if currentPage > 0 {
				currentPage--
				updatePage()
			}
			return nil
		case tcell.KeyRight:
			if currentPage < maxPages {
				currentPage++
				updatePage()
			}
			return nil
		}
		return event
	})

	if err := app.SetRoot(mainFlex, true).Run(); err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}
}
