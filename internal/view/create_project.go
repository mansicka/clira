package view

import (
	"log"

	"github.com/mansicka/rtpms/internal/project"
	"github.com/mansicka/rtpms/internal/ui"
	"github.com/mansicka/rtpms/internal/view/modal"
	"github.com/rivo/tview"
)

func ShowCreateProjectForm(ui *ui.UIManager) {
	form := tview.NewForm().
		AddInputField("Project Name", "", 50, nil, nil).
		AddInputField("Project Key", "", 10, nil, nil).
		AddTextArea("Description", "", 50, 5, 0, nil).
		AddInputField("Client", "", 50, nil, nil)
	form.AddButton("Create", func() {
		name := form.GetFormItemByLabel("Project Name").(*tview.InputField).GetText()
		projectKey := form.GetFormItemByLabel("Project Key").(*tview.InputField).GetText()
		description := form.GetFormItemByLabel("Description").(*tview.TextArea).GetText()
		client := form.GetFormItemByLabel("Client").(*tview.InputField).GetText()

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
			Status:       "active",
			ActiveSprint: 0,
			Users:        make(map[string]string),
		}

		err := project.SaveProject(newProject)
		if err != nil {
			modal.ShowErrorModal(ui, "Failed to create project")
		}

		//ShowProjectList(app, pages)
	}).
		AddButton("Cancel", func() {
			ui.GoBack()
		})

	form.SetTitle(" Create New Project ").SetBorder(true)
	ui.AddView("create_project", form)
}
