package git

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/mansicka/rtpms/internal/storage"
)

func InitializeGitRepository() error {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return fmt.Errorf("failed to get storage instance: %w", err)
	}

	repositoryExists, err := RepositoryExists()
	if err != nil {
		return fmt.Errorf("error checking if repository exists: %w", err)
	}

	if !repositoryExists {
		_, err := git.PlainInit(storageInstance.RootDir, false)
		if err != nil {
			return fmt.Errorf("git init failed: %w", err)
		}
		fmt.Println("Initialized new Git repository")
	}

	if err := CreateGitIgnore(); err != nil {
		return fmt.Errorf("error creating .gitignore: %w", err)
	}

	hasChanges, err := HasUncommittedChanges()
	if err != nil {
		return fmt.Errorf("error checking git status: %w", err)
	}

	if hasChanges {
		err = DoCommit("Initialized Clira repository")
		if err != nil {
			return fmt.Errorf("error committing initial state: %w", err)
		}
	} else {
		fmt.Println("No changes to commit.")
	}

	return nil
}

func RepositoryExists() (bool, error) {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return false, fmt.Errorf("failed to get storage instance: %v", err)
	}

	repoPath := storageInstance.RootDir
	_, err = git.PlainOpen(repoPath)
	if err == git.ErrRepositoryNotExists {
		return false, nil
	}
	return err == nil, err
}

func RemoteExists() (bool, error) {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return false, fmt.Errorf("failed to get storage instance: %v", err)
	}

	repo, err := git.PlainOpen(storageInstance.RootDir)
	if err != nil {
		return false, fmt.Errorf("failed to open repository: %v", err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return false, fmt.Errorf("failed to get remotes: %v", err)
	}

	return len(remotes) > 0, nil
}

func DoCommit(message string) error {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return fmt.Errorf("failed to get storage instance: %v", err)
	}

	repo, err := git.PlainOpen(storageInstance.RootDir)
	if err != nil {
		return fmt.Errorf("failed to open repository: %v", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %v", err)
	}

	err = wt.AddGlob("*")
	if err != nil {
		return fmt.Errorf("failed to add files: %v", err)
	}

	_, err = wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Clira Auto Commit",
			Email: "auto@clira.git",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %v", err)
	}

	fmt.Println("commit successful:", message)
	return nil
}

func DoPush() error {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return fmt.Errorf("failed to get storage instance: %v", err)
	}

	repo, err := git.PlainOpen(storageInstance.RootDir)
	if err != nil {
		return fmt.Errorf("failed to open repository: %v", err)
	}

	remoteExists, err := RemoteExists()
	if err != nil || !remoteExists {
		return fmt.Errorf("no remote repository found")
	}

	// TODO: auth to git
	auth := &http.BasicAuth{
		Username: "your-github-username",
		Password: "your-personal-access-token",
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		return fmt.Errorf("failed to push changes: %v", err)
	}

	fmt.Println("Push successful!")
	return nil
}

func CreateGitIgnore() error {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return fmt.Errorf("failed to get storage instance: %v", err)
	}

	gitignorePath := filepath.Join(storageInstance.RootDir, ".gitignore")
	gitignoreContent := "clira*\n"

	if !storageInstance.FileExists(gitignorePath) {
		return storageInstance.WriteFile(".gitignore", []byte(gitignoreContent))
	}

	return nil
}

func DoCommitAndPush(message string) error {
	if err := DoCommit(message); err != nil {
		return err
	}
	return DoPush()
}

func HasUncommittedChanges() (bool, error) {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return false, fmt.Errorf("failed to get storage instance: %w", err)
	}

	repo, err := git.PlainOpen(storageInstance.RootDir)
	if err != nil {
		return false, fmt.Errorf("failed to open git repository: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get git worktree: %w", err)
	}

	status, err := wt.Status()
	if err != nil {
		return false, fmt.Errorf("failed to get git status: %w", err)
	}

	return !status.IsClean(), nil
}