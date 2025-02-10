package organization

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/mansicka/clira/internal/storage"
)

// Organization struct holds the organization's data
type Organization struct {
	Name        string   `json:"name"`
	Created     string   `json:"created"`
	Description string   `json:"description"`
	Admins      []string `json:"admins"`
}

var orgFile string = "organization.json"

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// LoadOrganization safely loads organization data
func LoadOrganization() (*Organization, error) {

	storageInstance, err := storage.GetStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage instance: %w", err)
	}

	// Check if the file exists
	if storageInstance.FileExists(orgFile) {
		data, err := storageInstance.ReadFile(orgFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read organization.json: %w", err)
		}

		// If the file is empty, return an error
		if len(data) == 0 {
			return nil, errors.New("organization.json is empty")
		}

		var org Organization
		if err := json.Unmarshal(data, &org); err != nil {
			return nil, fmt.Errorf("failed to parse organization.json: %w", err)
		}

		// Extra safety check: Ensure Admins field is always initialized
		if org.Admins == nil {
			org.Admins = []string{}
		}

		return &org, nil
	}
	return nil, errors.New("organization.json not found")
}

// SaveOrganization creates a new organization and writes it to JSON file
func SaveOrganization(name, description string) error {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return fmt.Errorf("failed to get storage instance: %w", err)
	}

	org := Organization{
		Name:        name,
		Created:     time.Now().UTC().Format(time.RFC3339),
		Description: description,
		Admins:      []string{}, // Admin roles to be added after
	}

	data, err := json.MarshalIndent(org, "", "  ")
	if err != nil {
		return err
	}

	filewriteerr := storageInstance.WriteFile(orgFile, data)
	if filewriteerr != nil {
		return filewriteerr
	}
	//return storage.InitializeGitRepository()
	return nil

}

func AddAdmin(username string) error {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return fmt.Errorf("failed to get storage instance: %w", err)
	}

	data, err := storageInstance.ReadFile(orgFile)
	if err != nil {
		return err
	}

	var org Organization
	if err := json.Unmarshal(data, &org); err != nil {
		return err
	}

	// Check if user exists in users/
	userFile := fmt.Sprintf("users/%s.json", username)
	if !storageInstance.FileExists(userFile) {
		return errors.New(fmt.Sprintf("user '%s' does not exist", userFile))
	}

	// Check if already an admin
	for _, admin := range org.Admins {
		if admin == username {
			return errors.New(fmt.Sprintf("user '%s' is already an admin", username))
		}
	}

	// Add new admin
	org.Admins = append(org.Admins, username)

	// Save updated organization data
	updatedData, err := json.MarshalIndent(org, "", "  ")
	if err != nil {
		return err
	}

	err = storageInstance.WriteFile(orgFile, updatedData)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' added as an admin.", username)
	return nil
}
