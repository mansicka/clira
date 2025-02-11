package view

import (
	"fmt"
	"log"

	"github.com/mansicka/rtpms/internal/project"
	"github.com/mansicka/rtpms/internal/ui"
	"github.com/mansicka/rtpms/internal/util"
	"github.com/mansicka/rtpms/internal/view/modal"
	"github.com/rivo/tview"
)

func ShowEditProjectForm(ui *ui.UIManager, proj project.Project) {
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
			modal.ShowErrorModal(ui, fmt.Sprintf("Failed to save project: %s", err))
		}

		InitProjectList(ui)
	}).
		AddButton("Cancel", func() {
			InitProjectList(ui)
		})

	form.SetTitle(fmt.Sprintf("Edit project: %s - %s", proj.ProjectKey, proj.Name)).SetBorder(true)
	ui.AddView("edit_project", form)
}
