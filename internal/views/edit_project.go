package views

import (
	"fmt"
	"log"

	"github.com/mansicka/rtpms/internal/project"
	"github.com/mansicka/rtpms/internal/ui"
	"github.com/mansicka/rtpms/internal/util"
	"github.com/rivo/tview"
)

func ShowEditProjectForm(app *tview.Application, pages *tview.Pages, proj project.Project) {
	form := tview.NewForm().
		AddInputField("Project Name", proj.Name, 50, nil, nil).
		AddInputField("Project Key", proj.ProjectKey, 10, nil, nil).
		AddTextArea("Description", proj.Description, 50, 5, 0, nil).
		AddInputField("Client", proj.Client, 50, nil, nil)
	options := []string{"active", "inactive"}
	form.AddDropDown("Status", options, util.FindIndex(options, proj.Status), nil)
	form.GetFormItem(1).SetDisabled(true)

	form.AddButton("Save", func() {
		name := form.GetFormItemByLabel("Project Name").(*tview.InputField).GetText()
		projectKey := form.GetFormItemByLabel("Project Key").(*tview.InputField).GetText()
		description := form.GetFormItemByLabel("Description").(*tview.TextArea).GetText()
		client := form.GetFormItemByLabel("Client").(*tview.InputField).GetText()
		_, status := form.GetFormItemByLabel("Status").(*tview.DropDown).GetCurrentOption()

		if name == "" || projectKey == "" {
			log.Println("Project name and key are required!")
			return
		}

		if client == "" {
			client = "internal"
		}

		newProject := project.Project{
			ID:           projectKey,
			ProjectKey:   projectKey,
			Name:         name,
			Description:  description,
			Client:       client,
			Status:       status,
			ActiveSprint: proj.ActiveSprint,
			Users:        proj.Users,
		}

		err := project.EditProject(newProject)
		if err != nil {
			ui.ShowErrorModal(app, pages, fmt.Sprintf("Failed to save project: %s", err), "create_project")
		}

		ShowProjectList(app, pages)
	}).
		AddButton("Cancel", func() {
			ShowProjectList(app, pages)
		})

	form.SetTitle(fmt.Sprintf("Edit project: %s - %s", proj.ProjectKey, proj.Name)).SetBorder(true)
	pages.AddPage("edit_project", form, true, true)
	app.SetRoot(form, true)
}
