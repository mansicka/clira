package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mansicka/rtpms/internal/event"
	"github.com/mansicka/rtpms/internal/git"
	"github.com/mansicka/rtpms/internal/organization"
	"github.com/mansicka/rtpms/internal/storage"
	"github.com/mansicka/rtpms/internal/ui"
	"github.com/mansicka/rtpms/internal/view"
)

func main() {
	_ = godotenv.Load()

	rootDir := os.Getenv("RTPMS_ROOTDIR")

	if rootDir == "" {
		panic("ss")
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

	uiManager := ui.NewUIManager()

	view.InitCreateOrganizationForm(uiManager)
	view.InitCreateAdminUserForm(uiManager)
	view.ShowLoginPrompt(uiManager)

	orgData, err := organization.LoadOrganization()

	if err != nil {
		log.Print(err)
	}
	if orgData == nil {
		uiManager.SwitchToView("create_organization")
	} else {
		if len(orgData.Admins) == 0 {
			uiManager.SwitchToView("create_admin_user")
		} else {
			uiManager.SwitchToView("login")
		}
	}

	if err := uiManager.Run(); err != nil {
		log.Fatal(err)
	}
}
