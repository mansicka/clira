package project

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mansicka/ugh/internal/storage"
)

type Project struct {
	ID           string            `json:"id"`
	ProjectKey   string            `json:"project_key"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Client       string            `json:"client"`
	Status       string            `json:"status"`
	ActiveSprint int               `json:"active_sprint"`
	Users        map[string]string `json:"users"`
}

func GetAllProjects() ([]Project, error) {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage instance: %w", err)
	}

	projectDirs, err := storageInstance.ReadDir("projects")
	if err != nil {
		return nil, fmt.Errorf("failed to read projects directory: %w", err)
	}

	var projects []Project

	for _, dir := range projectDirs {
		if dir.IsDir() {
			project, err := GetProject(dir.Name())
			if err != nil {
				fmt.Errorf(err.Error())
				continue
			}
			projects = append(projects, project)
		}
	}

	return projects, nil
}

func SaveProject(project Project) error {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return fmt.Errorf("failed to get storage instance: %v", err)
	}

	projectDir := fmt.Sprintf("projects/%s", project.ProjectKey)
	err = storageInstance.CreateDir(projectDir)
	if err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	projectFilePath := fmt.Sprintf("%s/project.json", projectDir)
	projectData, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project data: %v", err)
	}

	err = storageInstance.WriteFile(projectFilePath, projectData)
	if err != nil {
		return fmt.Errorf("failed to save project: %v", err)
	}
	return nil
}

func GetProject(projectKey string) (Project, error) {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return Project{}, fmt.Errorf("failed to get storage instance: %v", err)
	}

	projectFilePath := fmt.Sprintf("projects/%s/project.json", projectKey)
	if storageInstance.FileExists(projectFilePath) {
		projectData, err := storageInstance.ReadFile(projectFilePath)
		if err != nil {
			return Project{}, fmt.Errorf("could not read project file %s: %v", projectFilePath, err)
		}
		var proj Project
		err = json.Unmarshal(projectData, &proj)
		if err != nil {
			return Project{}, fmt.Errorf("failed to parse project.json for %s: %v", projectFilePath, err)
		}
		return proj, nil
	}
	return Project{}, fmt.Errorf("could not read project file %s: %v", projectFilePath, err)
}

func EditProject(project Project) error {
	existingProject, err := GetProject(project.ProjectKey)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(project, existingProject) {
		return fmt.Errorf("no changes detected")
	}

	saveErr := SaveProject(project)
	if saveErr != nil {
		return err
	}
	return nil
}
