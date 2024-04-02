package wbtask

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tomas2707/wbtask/repository"
	"log/slog"
	"net/http"
	"net/mail"
	"time"
)

// Service provides an api contracts to user-related operations.
type Service struct {
	log      *slog.Logger
	userRepo repository.UserManager
}

// NewService creates a new instance of the user Service.
func NewService(log *slog.Logger, userRepo repository.UserManager) *Service {
	return &Service{
		log:      log,
		userRepo: userRepo,
	}
}

// SaveUserHandler is an HTTP handler for saving a new user to the database.
// It reads the user details from the request body and attempts to create a new user record.
// If successful, it responds with the created user data.
// It responds with an appropriate HTTP error status in case of failure.
func (s *Service) SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.log.Error("Error decoding user data", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = uuid.Parse(user.ID)
	if err != nil {
		s.log.Error("Invalid user id", "id", user.ID, "error", err)
		http.Error(w, "Invalid user id: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = mail.ParseAddress(user.Email)
	if err != nil {
		s.log.Error("Invalid email address", "email", user.Email, "error", err)
		http.Error(w, "Invalid email address: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = time.Parse(time.RFC3339, user.DateOfBirth)
	if err != nil {
		s.log.Error("Invalid date format", "date", user.DateOfBirth, "error", err)
		http.Error(w, "Invalid date format: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = s.userRepo.Create(user)
	if err != nil {
		s.log.Error("Error creating user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.log.Info("User saved successfully", "userID", user.ID)
	json.NewEncoder(w).Encode(user)
}

// GetUserHandler is an HTTP handler for retrieving a user by their ID from the database.
// It extracts the user ID from the URL path, looks up the user in the database,
// and responds with the user data if found.
// It responds with an appropriate HTTP error status in case the user is not found or in case of other errors.
func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := s.userRepo.Get(id)
	if err != nil {
		s.log.Error("User not found", "error", err, "userID", id)
		http.NotFound(w, r)
		return
	}

	s.log.Info("User retrieved successfully", "userID", user.ID)
	json.NewEncoder(w).Encode(user)
}
