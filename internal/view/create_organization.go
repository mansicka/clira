package view

import (
	"fmt"
	"strings"

	"github.com/mansicka/rtpms/globals"
	"github.com/mansicka/rtpms/internal/organization"
	"github.com/mansicka/rtpms/internal/ui"
	"github.com/mansicka/rtpms/internal/view/modal"
	"github.com/rivo/tview"
)

// InitCreateOrganizationForm displays the organization creation UI
func InitCreateOrganizationForm(ui *ui.UIManager) {

	header := tview.NewTextView().
		SetText(fmt.Sprintf("%s", globals.LogoAndHeaderText+"\n"+globals.CreateOrganizationInfo)).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	form := tview.NewForm().
		AddInputField("Organization Name:", "", 50, nil, nil).
		AddInputField("Description:", "", 50, nil, nil)
	form.AddButton("Create", func() {
		name := strings.TrimSpace(form.GetFormItemByLabel("Organization Name:").(*tview.InputField).GetText())
		description := strings.TrimSpace(form.GetFormItemByLabel("Description:").(*tview.InputField).GetText())

		if name == "" {
			modal.ShowErrorModal(ui, "Organization name is required!")
		}

		err := organization.SaveOrganization(name, description)
		if err != nil {
			modal.ShowErrorModal(ui, fmt.Sprintf("Error saving organization: %s", err))
		}

		//ShowCreateAdminUserForm(app, pages)
		ui.SwitchToView("create_admin_user")
	}).
		AddButton("Cancel", func() {
			modal.ShowExitConfirmationModal(ui)
		}).
		SetTitle("Create Organization").
		SetBorder(true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 16, 1, false).
		AddItem(form, 0, 1, true)

	ui.AddView("create_organization", layout)
	ui.SwitchToView("create_organization")
}
