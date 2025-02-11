package state

import (
	"sync"

	"github.com/mansicka/rtpms/internal/project"
	"github.com/mansicka/rtpms/internal/user"
)

// AppState keeps track of global application state
type AppState struct {
	mu              sync.RWMutex
	CurrentUser     *user.User
	SelectedProject *project.Project
}

var appState *AppState
var once sync.Once

// GetState returns a singleton instance of the AppState
func GetState() *AppState {
	once.Do(func() {
		appState = &AppState{
			CurrentUser:     nil,
			SelectedProject: nil,
		}
	})
	return appState
}

// SetUser updates the logged-in user
func (s *AppState) SetUser(user *user.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CurrentUser = user
}

// SetProject updates the selected project
func (s *AppState) SetProject(proj *project.Project) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SelectedProject = proj
}

// GetUser returns the logged-in user
func (s *AppState) GetUser() *user.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.CurrentUser
}

// GetProject returns the selected project
func (s *AppState) GetProject() *project.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.SelectedProject
}
