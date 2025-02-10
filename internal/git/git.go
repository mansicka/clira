package git

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/mansicka/ugh/internal/storage"
)

func InitializeGitRepository() error {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return fmt.Errorf("failed to get storage instance: %v", err)
	}

	repoPath := storageInstance.RootDir

	if exists, _ := RepositoryExists(); exists {
		fmt.Println("Git repository already exists.")
		return nil
	}

	_, err = git.PlainInit(repoPath, false)
	if err != nil {
		return fmt.Errorf("git init failed: %w", err)
	}

	fmt.Println("Initialized new Git repository in", repoPath)

	if err := CreateGitIgnore(); err != nil {
		return fmt.Errorf("error creating .gitignore: %w", err)
	}

	err = DoCommit("Initialized UGH repository")
	if err != nil {
		return fmt.Errorf("error committing initial state: %w", err)
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
			Name:  "UGH Auto Commit",
			Email: "auto@ugh.git",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %v", err)
	}

	fmt.Println("✅ Commit successful:", message)
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
	gitignoreContent := "ugh*\n"

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
