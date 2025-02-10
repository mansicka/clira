package views

import (
	"log"

	"github.com/mansicka/ugh/internal/project"
	"github.com/mansicka/ugh/internal/ui"
	"github.com/rivo/tview"
)

func ShowCreateProjectForm(app *tview.Application, pages *tview.Pages) {
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
			ui.ShowErrorModal(app, pages, "Failed to create project", "create_project")
		}

		ShowProjectList(app, pages)
	}).
		AddButton("Cancel", func() {
			ShowProjectList(app, pages)
		})

	form.SetTitle(" Create New Project ").SetBorder(true)
	pages.AddPage("create_project", form, true, true)
	app.SetRoot(form, true)
}
