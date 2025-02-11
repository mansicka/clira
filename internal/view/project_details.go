package view

import (
	"fmt"

	"github.com/mansicka/rtpms/internal/project"
	"github.com/mansicka/rtpms/internal/ui"

	//"github.com/mansicka/rtpms/internal/state"
	"github.com/rivo/tview"
)

func InitProjectDetails(ui *ui.UIManager, proj *project.Project) {
	projectInfo := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf(
			"[green]Project Name:[white] %s\n[green]Project Key:[white] %s\n[green]Description:[white] %s\n[green]Client:[white] %s\n[green]Status:[white] %s\n[green]Active Sprint:[white] %d",
			proj.Name, proj.ProjectKey, proj.Description, proj.Client, proj.Status, proj.ActiveSprint,
		)).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	menu := tview.NewList().
		AddItem("Return to Project List", "", 'r', func() {
			InitProjectList(ui)
		}).
		AddItem("Exit", "", 'q', func() {
			ui.App.Stop()
		})

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(projectInfo, 0, 1, false).
		AddItem(menu, 0, 1, true)

	ui.AddView("project_details", layout)
}
