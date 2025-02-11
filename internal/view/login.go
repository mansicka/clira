package view

import (
	"fmt"

	"github.com/mansicka/rtpms/globals"
	"github.com/mansicka/rtpms/internal/organization"
	"github.com/mansicka/rtpms/internal/state"
	"github.com/mansicka/rtpms/internal/ui"
	"github.com/mansicka/rtpms/internal/user"
	"github.com/mansicka/rtpms/internal/view/modal"
	"github.com/rivo/tview"
)

// ShowLoginPrompt displays a login prompt if no user is logged in
func ShowLoginPrompt(ui *ui.UIManager) {
	organization, err := organization.LoadOrganization()
	if err != nil {
		modal.ShowErrorModal(ui, "Getting organization data failed. Probably corrupted json. Good luck!")
		return
	}

	header := tview.NewTextView().
		SetText(fmt.Sprintf(globals.LogoAndHeaderText+"Organization: %s", organization.Name)).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	form := tview.NewForm().
		AddInputField("Username:", "", 30, nil, nil).
		AddPasswordField("Password:", "", 30, '*', nil)
	form.AddButton("Login", func() {
		username := form.GetFormItemByLabel("Username:").(*tview.InputField).GetText()
		password := form.GetFormItemByLabel("Password:").(*tview.InputField).GetText()

		loggedInUser, err := user.LoginUser(username, password)
		if err != nil {
			modal.ShowErrorModal(ui, fmt.Sprintf("Login failed: %s", err))
			return
		}

		appState := state.GetState()
		appState.SetUser(loggedInUser)

		InitMainMenu(ui)
		ui.SwitchToView("main_menu")
	}).
		AddButton("Exit", func() {
			modal.ShowExitConfirmationModal(ui)
		}).
		SetTitle("User Login").
		SetBorder(true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 16, 1, false).
		AddItem(form, 0, 1, true)

	ui.AddView("login", layout)
}
