package repository

// UserManager defines the interface for managing user entities.
type UserManager interface {
	Create(user User) error
	Get(id string) (User, error)
}

// User represents a user entity in the system.
// It includes fields for the user's ID, name, email, and date of birth,
// each tagged for JSON serialization and GORM modeling.
type User struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Email       string `json:"email" gorm:"unique"`
	DateOfBirth string `json:"date_of_birth"`
}
