package view

import (
	"fmt"

	"github.com/mansicka/rtpms/internal/project"
	"github.com/mansicka/rtpms/internal/state"
	"github.com/mansicka/rtpms/internal/ui"
	"github.com/mansicka/rtpms/internal/view/modal"
	"github.com/rivo/tview"
)

func InitProjectList(uiManager *ui.UIManager) {
	appState := state.GetState()

	projectList := tview.NewList().ShowSecondaryText(false)

	projectDetails := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Select a project from the list")

	projects, err := project.GetAllProjects()
	if err != nil {
		modal.ShowErrorModal(uiManager, "Failed to get projects")
		return
	}

	if len(projects) != 0 {
		for _, proj := range projects {
			itemText := fmt.Sprintf(" 	[ %s ] : %s", proj.ProjectKey, proj.Name)
			projectList.AddItem(itemText, "", 0, nil)
		}

		projectList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			if index >= 0 && index < len(projects) {
				updateProjectDetails(projectDetails, &projects[index])
			}
		})

		projectList.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			if index >= 0 && index < len(projects) {
				appState.SetProject(&projects[index])
				ShowEditProjectForm(uiManager, projects[index])
			}
		})
	} else {
		projectList.AddItem("No projects found", "", 0, nil)
	}

	projectList.SetBorder(true).SetTitle(" Projects ")

	createProjectButton := tview.NewButton(" Create New Project ").
		SetSelectedFunc(func() {
			uiManager.SwitchToView("create_project")
		})

	backButton := tview.NewButton(" Back to Main Menu ").
		SetSelectedFunc(func() {
			uiManager.SwitchToView("main_menu")
		})

	footer := tview.NewTextView().
		SetText("[Arrow Keys] Select Project  [Enter] Edit Project  [ESC] Back").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	buttons := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(createProjectButton, 20, 1, false).
		AddItem(backButton, 20, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(projectList, 30, 1, true).
				AddItem(projectDetails, 0, 2, false),
			0, 1, true).
		AddItem(buttons, 3, 1, false).
		AddItem(footer, 1, 1, false)

	uiManager.AddView("project_list", layout)
}

func updateProjectDetails(detailsView *tview.TextView, proj *project.Project) {
	detailsView.SetText(fmt.Sprintf(
		"[green]Project Name:[white] %s\n"+
			"[green]Project Key:[white] %s\n"+
			"[green]Description:[white] %s\n"+
			"[green]Client:[white] %s\n"+
			"[green]Status:[white] %s\n"+
			"[green]Active Sprint:[white] %d",
		proj.Name, proj.ProjectKey, proj.Description, proj.Client, proj.Status, proj.ActiveSprint,
	))
	detailsView.SetBorder(true).SetTitle(" Project information ")
}
