package gorm

import (
	"fmt"
	"github.com/tomas2707/wbtask/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConfig holds the database configuration parameters
type DBConfig struct {
	Host     string
	User     string
	DBName   string
	SSLMode  string
	Password string
}

// InitDB initializes and returns a database connection.
func InitDB(cfg DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host,
		cfg.User,
		cfg.DBName,
		cfg.SSLMode,
		cfg.Password,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&repository.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
