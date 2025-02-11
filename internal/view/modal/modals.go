package modal

import "github.com/mansicka/rtpms/internal/ui"

func ShowErrorModal(uiManager *ui.UIManager, message string) {
	modal := NewDynamicModal(uiManager, message, []string{"OK"}, []func(){uiManager.GoBack})
	uiManager.AddView("error_modal", modal)
	uiManager.SwitchToView("error_modal")
}

func ShowExitConfirmationModal(uiManager *ui.UIManager) {
	modal := NewDynamicModal(uiManager,
		"Are you sure you want to exit computer program?",
		[]string{"Yes", "No"}, []func(){uiManager.App.Stop, uiManager.GoBack})
	uiManager.AddView("exit_confirmation_modal", modal)
	uiManager.SwitchToView("exit_confirmation_modal")
}
