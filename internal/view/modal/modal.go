package modal

import (
	"github.com/mansicka/rtpms/internal/ui"
	"github.com/rivo/tview"
)

func NewDynamicModal(ui *ui.UIManager, message string, buttonLabels []string, actions []func()) tview.Primitive {
	modal := tview.NewModal().
		SetText(message)

	for _, btn := range buttonLabels {
		modal.AddButtons([]string{btn})
	}

	modal.SetDoneFunc(func(bIndex int, buttonLabel string) {
		if bIndex == -1 {
			ui.GoBack()
		}
		if bIndex < len(actions) && actions[bIndex] != nil {
			actions[bIndex]()
		}
	})

	modalContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 2, false).
		AddItem(modal, 10, 1, true).
		AddItem(nil, 0, 2, false)

	ui.App.SetFocus(modal)
	return modalContainer
}
