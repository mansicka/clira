package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewApp creates a new tview application with page navigation
func NewApp() (*tview.Application, *tview.Pages) {
	var Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorNone,
		ContrastBackgroundColor:     tcell.ColorBlue,
		MoreContrastBackgroundColor: tcell.ColorGreen,
		BorderColor:                 tcell.ColorWhite,
		TitleColor:                  tcell.ColorWhite,
		GraphicsColor:               tcell.ColorWhite,
		PrimaryTextColor:            tcell.ColorWhite,
		SecondaryTextColor:          tcell.ColorYellow,
		TertiaryTextColor:           tcell.ColorGreen,
		InverseTextColor:            tcell.ColorBlue,
		ContrastSecondaryTextColor:  tcell.ColorNavy,
	}
	tview.Styles = Styles
	app := tview.NewApplication()
	pages := tview.NewPages()

	var history []string

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			if len(history) > 1 {
				history = history[:len(history)-1]
				pages.SwitchToPage(history[len(history)-1])
				return nil
			}
			app.Stop()
			return nil
		}
		return event
	})

	return app, pages
}

func ShowErrorModal(app *tview.Application, pages *tview.Pages, message string, returnPage string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.SwitchToPage(returnPage)
		})
	pages.AddPage("error_modal", modal, true, true)
}
