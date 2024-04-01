package gorm

import (
	"errors"
	"github.com/tomas2707/wbtask/repository"
	"gorm.io/gorm"
)

// UserRepository represents a repository for managing Users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create adds a new User to the database
func (repo *UserRepository) Create(user repository.User) error {
	return repo.db.Create(&user).Error
}

// Get retrieves a User by ID
// If the user with the specified ID does not exist, it returns gorm.ErrRecordNotFound.
func (repo *UserRepository) Get(id string) (repository.User, error) {
	var user repository.User
	err := repo.db.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.User{}, err
	}

	return user, nil
}
