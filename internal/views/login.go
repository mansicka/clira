package views

import (
	"fmt"

	"github.com/mansicka/ugh/globals"
	"github.com/mansicka/ugh/internal/organization"
	"github.com/mansicka/ugh/internal/state"
	"github.com/mansicka/ugh/internal/ui"
	"github.com/mansicka/ugh/internal/user"
	"github.com/rivo/tview"
)

// ShowLoginPrompt displays a login prompt if no user is logged in
func ShowLoginPrompt(app *tview.Application, pages *tview.Pages) {
	organization, err := organization.LoadOrganization()
	if err != nil {
		ui.ShowErrorModal(app, pages, "Getting organization data failed. Probably corrupted json. Good luck!",
			"login")
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
			errorModal := tview.NewModal().
				SetText(fmt.Sprintf("Login failed: %s", err)).
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetRoot(form, true).SetFocus(form)
				})
			app.SetRoot(errorModal, true).SetFocus(errorModal)
			return
		}

		appState := state.GetState()
		appState.SetUser(loggedInUser)

		ShowMainMenu(app, pages)
		pages.SwitchToPage("main_menu")
	}).
		AddButton("Exit", func() {
			app.Stop()
		}).
		SetTitle("User Login").
		SetBorder(true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 16, 1, false).
		AddItem(form, 0, 1, true)

	pages.AddPage("login", layout, true, false)
}
