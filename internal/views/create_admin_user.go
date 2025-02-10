package views

import (
	"fmt"

	"github.com/mansicka/clira/globals"
	"github.com/mansicka/clira/internal/organization"
	"github.com/mansicka/clira/internal/ui"
	"github.com/mansicka/clira/internal/user"
	"github.com/rivo/tview"
)

// ShowCreateAdminUserForm displays the admin user creation UI
func ShowCreateAdminUserForm(app *tview.Application, pages *tview.Pages) {
	header := tview.NewTextView().
		SetText(fmt.Sprintf("%s", globals.LogoAndHeaderText+"\n"+globals.CreateAdminUserInfo)).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	form := tview.NewForm().
		AddInputField("Admin Username", "", 50, nil, nil).
		AddInputField("Password", "", 50, nil, nil).
		AddInputField("Password again", "", 50, nil, nil)

	form.AddButton("Create", func() {
		usernameField, _ := form.GetFormItem(0).(*tview.InputField)
		passwordField, _ := form.GetFormItem(1).(*tview.InputField)
		passwordAgainField, _ := form.GetFormItem(2).(*tview.InputField)

		if usernameField == nil || passwordField == nil {
			ui.ShowErrorModal(app, pages, "Error: Failed to retrieve form inputs", "create_admin_user")
			return
		}

		username := usernameField.GetText()
		password := passwordField.GetText()
		passwordAgain := passwordAgainField.GetText()

		if password != passwordAgain {
			ui.ShowErrorModal(app, pages, "Passwords do not match!", "create_admin_user")
			return
		}

		if username == "" || password == "" {
			ui.ShowErrorModal(app, pages, "Username and Password are required!", "create_admin_user")
			return
		}

		if err := user.SaveUser(username, password, "admin"); err != nil {
			ui.ShowErrorModal(app, pages, fmt.Sprintf("Error saving user:\n%s", err), "create_admin_user")
			return
		}

		if err := organization.AddAdmin(username); err != nil {
			ui.ShowErrorModal(app, pages, fmt.Sprintf("Error adding admin:\n%s", err), "create_admin_user")
			return
		}
		ShowLoginPrompt(app, pages)
		pages.SwitchToPage("login")
	}).
		SetTitle("Create Admin User").
		SetBorder(true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 16, 1, false).
		AddItem(form, 0, 1, true)

	pages.AddPage("create_admin_user", layout, true, false)
}
