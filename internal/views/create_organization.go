package views

import (
	"fmt"
	"strings"

	"github.com/mansicka/ugh/globals"
	"github.com/mansicka/ugh/internal/organization"
	"github.com/mansicka/ugh/internal/ui"
	"github.com/rivo/tview"
)

// ShowCreateOrganizationForm displays the organization creation UI
func ShowCreateOrganizationForm(app *tview.Application, pages *tview.Pages) {

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
			ui.ShowErrorModal(app, pages, "Organization name is required!", "create_organization")

		}

		err := organization.SaveOrganization(name, description)
		if err != nil {
			ui.ShowErrorModal(app, pages, "Error saving organization", "create_organization")
		}

		ShowCreateAdminUserForm(app, pages)
		pages.SwitchToPage("create_admin_user")
	}).
		AddButton("Cancel", func() {
			pages.SwitchToPage("main_menu")
		}).
		SetTitle("Create Organization").
		SetBorder(true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 16, 1, false).
		AddItem(form, 0, 1, true)

	pages.AddPage("create_organization", layout, true, false)
}
