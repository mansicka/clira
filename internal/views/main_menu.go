package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mansicka/rtpms/globals"
	"github.com/mansicka/rtpms/internal/organization"
	"github.com/mansicka/rtpms/internal/state"
	"github.com/rivo/tview"
)

func ShowMainMenu(app *tview.Application, pages *tview.Pages) {
	appState := state.GetState()
	user := appState.GetUser()
	projectPtr := appState.GetProject()
	selectedProject := "No project selected"
	if projectPtr != nil {
		selectedProject = projectPtr.Name
	}

	orgData, err := organization.LoadOrganization()
	organizationName := "Unknown Organization"
	if err == nil && orgData != nil {
		organizationName = orgData.Name
	}

	header := tview.NewTextView().
		SetText(fmt.Sprintf(globals.Logo+"\n"+
			"👤 User: %s / %s\n"+
			"🏢 Organization: %s\n"+
			"📁 Project: %s\n"+
			"-----------------------------------",
			user.Username, user.Role, organizationName, selectedProject,
		)).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	menu := tview.NewList().
		AddItem("🗂️  Projects", "Manage projects", 'p', func() {
			ShowProjectList(app, pages)
		}).
		AddItem("🔃  Sprints", "Manage sprints", 's', func() {
		}).
		AddItem("📄  Issues", "Manage tickets", 'i', func() {
		}).
		AddItem("👨‍👨‍👦‍👦  Users", "Manage users", 'u', func() {
		}).
		AddItem("🛠️  Configuration", "Configure application", 'c', func() {
		}).
		AddItem("❌  Exit", "Quit application", 'q', func() {
			app.Stop()
		})

	menu.SetSelectedBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	footer := tview.NewTextView().
		SetText(globals.FooterNavigationInfo).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 12, 1, false).
		AddItem(menu, 0, 1, true).
		AddItem(footer, 1, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.Stop()
			return nil
		}
		return event
	})

	pages.AddPage("main_menu", layout, true, false)
}
