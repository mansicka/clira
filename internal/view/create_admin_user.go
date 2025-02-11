package view

import (
	"fmt"

	"github.com/mansicka/rtpms/globals"
	"github.com/mansicka/rtpms/internal/organization"
	"github.com/mansicka/rtpms/internal/ui"
	"github.com/mansicka/rtpms/internal/user"
	"github.com/mansicka/rtpms/internal/view/modal"
	"github.com/rivo/tview"
)

// InitCreateAdminUserForm displays the admin user creation UI
func InitCreateAdminUserForm(uiManager *ui.UIManager) {
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
			modal.ShowErrorModal(uiManager, "Error: Failed to retrieve form inputs")
			return
		}

		username := usernameField.GetText()
		password := passwordField.GetText()
		passwordAgain := passwordAgainField.GetText()

		if password != passwordAgain {
			modal.ShowErrorModal(uiManager, "Passwords don't match!")
			return
		}

		if username == "" || password == "" {
			modal.ShowErrorModal(uiManager, "Username and Password are both required.")
			return
		}

		if err := user.SaveUser(username, password, "admin"); err != nil {
			modal.ShowErrorModal(uiManager, fmt.Sprintf("Error saving user:\n%s", err))
			return
		}

		if err := organization.AddAdmin(username); err != nil {
			modal.ShowErrorModal(uiManager, fmt.Sprintf("Error adding admin:\n%s", err))
			return
		}

		ShowLoginPrompt(uiManager)
		uiManager.SwitchToView("login")
	}).
		SetTitle("Create Admin User").
		SetBorder(true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 16, 1, false).
		AddItem(form, 0, 1, true)

	uiManager.AddView("create_admin_user", layout)
}
