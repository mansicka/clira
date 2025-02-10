package user

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mansicka/ugh/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var userPath string = "users/"

// User struct defines a user in the system
type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

// HashPassword securely hashes the given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// SaveUser creates a new user JSON file with a hashed password
func SaveUser(username, password, role string) error {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return fmt.Errorf("failed to get storage instance: %w", err)
	}

	// Hash the password before storing
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := User{
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	// Serialize user data
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}

	// Save user file
	userFilePath := fmt.Sprintf("%s%s.json", userPath, username)
	return storageInstance.WriteFile(userFilePath, data)
}

func LoadUser(username string) (*User, error) {
	storageInstance, err := storage.GetStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage instance: %w", err)
	}

	userFilePath := fmt.Sprintf("users/%s.json", username)

	if !storageInstance.FileExists(userFilePath) {
		return nil, errors.New("user does not exist")
	}

	data, err := storageInstance.ReadFile(userFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read user file: %w", err)
	}

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}

	return &user, nil
}

func ValidateUser(username, password string) (bool, error) {
	user, err := LoadUser(username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("invalid password")
	}

	return true, nil
}

func GetUser(username string) (*User, error) {
	user, err := LoadUser(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func LoginUser(username, password string) (*User, error) {
	user, err := GetUser(username)
	if err != nil {
		return nil, fmt.Errorf("no user found for username %s", username)
	}

	valid, err := ValidateUser(username, password)
	if !valid {
		return nil, err
	}

	return user, nil
}
