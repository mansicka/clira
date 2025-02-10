package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UIManager struct {
	App     *tview.Application
	Pages   *tview.Pages
	Views   map[string]tview.Primitive
	History []string
}

func NewUIManager() *UIManager {
	app := tview.NewApplication()
	pages := tview.NewPages()

	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.NewRGBColor(43, 48, 59),
		ContrastBackgroundColor:     tcell.NewRGBColor(52, 61, 70),
		MoreContrastBackgroundColor: tcell.NewRGBColor(62, 71, 80),
		BorderColor:                 tcell.ColorWhite,
		TitleColor:                  tcell.ColorWhite,
		GraphicsColor:               tcell.ColorWhite,
		PrimaryTextColor:            tcell.ColorWhite,
		SecondaryTextColor:          tcell.NewRGBColor(186, 179, 117),
		TertiaryTextColor:           tcell.ColorGrey,
		InverseTextColor:            tcell.ColorBlack,
		ContrastSecondaryTextColor:  tcell.ColorBlack,
	}

	ui := &UIManager{
		App:     app,
		Pages:   pages,
		Views:   make(map[string]tview.Primitive),
		History: []string{},
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			ui.GoBack()
			return nil
		}
		return event
	})

	return ui
}

func (ui *UIManager) AddView(name string, view tview.Primitive) {
	ui.Views[name] = view
	ui.Pages.AddPage(name, view, true, false)
}

func (ui *UIManager) SwitchToView(name string) {
	if _, exists := ui.Views[name]; exists {
		if len(ui.History) == 0 || ui.History[len(ui.History)-1] != name {
			ui.History = append(ui.History, name)
		}
		ui.Pages.SwitchToPage(name)
	} else {
		log.Fatalf("View %s does not exist!", name)
	}
}

func (ui *UIManager) GoBack() {
	if len(ui.History) > 1 {
		ui.History = ui.History[:len(ui.History)-1]
		previousView := ui.History[len(ui.History)-1]
		ui.Pages.SwitchToPage(previousView)
	} else {
		ui.App.Stop()
	}
}

func (ui *UIManager) Run() error {
	return ui.App.SetRoot(ui.Pages, true).Run()
}
