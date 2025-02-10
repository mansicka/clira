package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mansicka/clira/internal/event"
	"github.com/mansicka/clira/internal/git"
	"github.com/mansicka/clira/internal/organization"
	"github.com/mansicka/clira/internal/state"
	"github.com/mansicka/clira/internal/storage"
	"github.com/mansicka/clira/internal/ui"
	"github.com/mansicka/clira/internal/views"
)

func main() {
	_ = godotenv.Load()

	rootDir := os.Getenv("CLIRA_ROOTDIR")

	if rootDir == "" {
		var err error
		rootDir, err = os.Getwd()
		if err != nil {
			log.Fatal("error getting current directory:", err)
		}
	}

	event.SetFileSaveEventListener(func(message string) {
		err := git.DoCommitAndPush(message)
		if err != nil {
			fmt.Errorf("event failed: %w", err)
		}
	})

	storageInstance, err := storage.NewStorage(rootDir)
	if err != nil {
		log.Panic(err)
	}

	if err := storageInstance.InitializeDirectoryStructure(); err != nil {
		log.Panic(err)
	}

	if err := git.InitializeGitRepository(); err != nil {
		log.Panic(err)
	}

	appState := state.GetState()
	user := appState.GetUser()

	app, pages := ui.NewApp()

	orgData, err := organization.LoadOrganization()
	if err != nil || orgData == nil {
		views.ShowCreateOrganizationForm(app, pages)
		pages.SwitchToPage("create_organization")
	} else {
		if len(orgData.Admins) == 0 {
			views.ShowCreateAdminUserForm(app, pages)
			pages.SwitchToPage("create_admin_user")
		} else if user == nil {
			views.ShowLoginPrompt(app, pages)
			pages.SwitchToPage("login")
		} else {
			views.ShowMainMenu(app, pages)
			pages.SwitchToPage("main_menu")
		}
	}

	app.SetRoot(pages, true)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
