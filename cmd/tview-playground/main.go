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

	// Header - modern styling with gradient-like colors
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	header.SetBackgroundColor(tcell.NewRGBColor(88, 86, 214)) // Modern purple gradient

	// Footer - modern styling
	footer := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("h/left: Previous | l/right: Next | q: Quit")
	footer.SetBackgroundColor(tcell.NewRGBColor(40, 42, 54)) // Darker, modern background

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
			// Page 0: Simple content with ROUNDED borders
			simpleBox := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("[yellow]Page 0:[white]\n\nSimple content box\n\n[gray]Notice the rounded corners!")
			simpleBox.SetBorder(true).
				SetBorderPadding(1, 1, 2, 2).
				SetBackgroundColor(tcell.NewRGBColor(68, 71, 90)).
				SetBorderColor(tcell.NewRGBColor(139, 233, 253)). // Cyan accent
				SetBorderStyle(tcell.StyleDefault)
			// Apply rounded border characters
			tview.Borders.Horizontal = '─'
			tview.Borders.Vertical = '│'
			tview.Borders.TopLeft = '╭'
			tview.Borders.TopRight = '╮'
			tview.Borders.BottomLeft = '╰'
			tview.Borders.BottomRight = '╯'
			contentArea.AddItem(simpleBox, 0, 1, false)

		case 1:
			// Page 1: Two horizontal panels with modern styling
			leftPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Left Panel\n\nFlex: 1")
			leftPanel.SetBorder(true).
				SetBorderPadding(1, 1, 2, 2).
				SetBackgroundColor(tcell.NewRGBColor(68, 71, 90)).
				SetBorderColor(tcell.NewRGBColor(255, 121, 198)) // Pink accent

			rightPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Right Panel\n\nFlex: 2")
			rightPanel.SetBorder(true).
				SetBorderPadding(1, 1, 2, 2).
				SetBackgroundColor(tcell.NewRGBColor(68, 71, 90)).
				SetBorderColor(tcell.NewRGBColor(139, 233, 253)) // Cyan accent

			contentArea.SetDirection(tview.FlexColumn).
				AddItem(leftPanel, 0, 1, false).
				AddItem(rightPanel, 0, 2, false)

		case 2:
			// Page 2: Three horizontal panels with gradient colors
			leftPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Left\n\nFlex: 1")
			leftPanel.SetBorder(true).
				SetBorderPadding(1, 1, 2, 2).
				SetBackgroundColor(tcell.NewRGBColor(68, 71, 90)).
				SetBorderColor(tcell.NewRGBColor(255, 121, 198)) // Pink

			centerPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Center\n\nFlex: 2")
			centerPanel.SetBorder(true).
				SetBorderPadding(1, 1, 2, 2).
				SetBackgroundColor(tcell.NewRGBColor(68, 71, 90)).
				SetBorderColor(tcell.NewRGBColor(189, 147, 249)) // Purple

			rightPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Right\n\nFlex: 1")
			rightPanel.SetBorder(true).
				SetBorderPadding(1, 1, 2, 2).
				SetBackgroundColor(tcell.NewRGBColor(68, 71, 90)).
				SetBorderColor(tcell.NewRGBColor(139, 233, 253)) // Cyan

			contentArea.SetDirection(tview.FlexColumn).
				AddItem(leftPanel, 0, 1, false).
				AddItem(centerPanel, 0, 2, false).
				AddItem(rightPanel, 0, 1, false)

		case 3:
			// Page 3: Vertical panels with modern colors
			topPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Top\nFlex: 1")
			topPanel.SetBorder(true).
				SetBorderPadding(0, 0, 2, 2).
				SetBackgroundColor(tcell.NewRGBColor(68, 71, 90)).
				SetBorderColor(tcell.NewRGBColor(80, 250, 123)) // Green

			middlePanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Middle\nFlex: 3")
			middlePanel.SetBorder(true).
				SetBorderPadding(0, 0, 2, 2).
				SetBackgroundColor(tcell.NewRGBColor(68, 71, 90)).
				SetBorderColor(tcell.NewRGBColor(241, 250, 140)) // Yellow

			bottomPanel := tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Bottom\nFlex: 1")
			bottomPanel.SetBorder(true).
				SetBorderPadding(0, 0, 2, 2).
				SetBackgroundColor(tcell.NewRGBColor(68, 71, 90)).
				SetBorderColor(tcell.NewRGBColor(255, 184, 108)) // Orange

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
