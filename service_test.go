package wbtask

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/tomas2707/wbtask/repository"
	repoGorm "github.com/tomas2707/wbtask/repository/gorm"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	db  *gorm.DB
	mux *http.ServeMux
}

func (s *ServiceTestSuite) SetupSuite() {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := repoGorm.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	var err error
	s.db, err = repoGorm.InitDB(cfg)
	s.NoError(err)

	userRepo := repoGorm.NewUserRepository(s.db)
	svc := NewService(logger, userRepo)
	s.mux = svc.Router()
}

func (s *ServiceTestSuite) SetupTest() {
	s.NoError(s.db.Migrator().DropTable(&repository.User{}))
	s.NoError(s.db.AutoMigrate(&repository.User{}))

	createCandidates := []repository.User{
		{
			ID:          "123e4567-e89b-12d3-a456-426614174000",
			Name:        "Alice Johnson",
			Email:       "alice.johnson@example.com",
			DateOfBirth: "1990-04-01T00:00:00Z",
		},
		{
			ID:          "c456e213-e89b-12d3-a456-426655440000",
			Name:        "Bob Smith",
			Email:       "bob.smith@example.com",
			DateOfBirth: "1985-05-23T00:00:00Z",
		},
	}
	for _, candidate := range createCandidates {
		s.NoError(s.db.Create(&candidate).Error)
	}
}

func (s *ServiceTestSuite) TestSaveUser() {
	candidates := []struct {
		description        string
		request            repository.User
		expectedErr        string
		expectedStatusCode int
	}{
		{
			description: "valid input",
			request: repository.User{
				ID:          "d789f456-e89b-12d3-a456-426614174000",
				Name:        "Charlie Davis",
				Email:       "charlie.davis@example.com",
				DateOfBirth: "1992-08-15T00:00:00Z",
			},
			expectedErr:        "",
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "duplicate uuid",
			request: repository.User{
				ID:          "123e4567-e89b-12d3-a456-426614174000",
				Name:        "Duplicate UUID",
				Email:       "duplicate.uuid@example.com",
				DateOfBirth: "1990-01-01T00:00:00Z",
			},
			expectedErr:        "duplicate key value violates unique constraint",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description: "invalid uuid format",
			request: repository.User{
				ID:          "invalid-uuid",
				Name:        "Invalid uuid",
				Email:       "invaliuuid",
				DateOfBirth: "1991-02-02T00:00:00Z",
			},
			expectedErr:        "invalid UUID length",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "invalid email format",
			request: repository.User{
				ID:          "284fab28-7586-412f-bef7-ca5d6ea596cf",
				Name:        "Invalid Email",
				Email:       "invalidemail",
				DateOfBirth: "1991-02-02T00:00:00Z",
			},
			expectedErr:        "mail: missing '@' or angle-addr",
			expectedStatusCode: http.StatusBadRequest,
		},
		// Candidate 4: email already exists
		{
			description: "email already exists",
			request: repository.User{
				ID:          "4456ab5b-c783-4ca8-b6d4-9ff8422bc2d0",
				Name:        "Email Exists",
				Email:       "alice.johnson@example.com",
				DateOfBirth: "1993-03-03T00:00:00Z",
			},
			expectedErr:        "duplicate key value violates unique constraint",
			expectedStatusCode: http.StatusInternalServerError,
		},
		// Candidate 5: invalid date format
		{
			description: "invalid date format",
			request: repository.User{
				ID:          "98c53896-7cff-401c-8ba0-77f6a0b96063",
				Name:        "Invalid Date",
				Email:       "invalid.date@example.com",
				DateOfBirth: "03-03-1993",
			},
			expectedErr:        "Invalid date format",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, candidate := range candidates {
		body, _ := json.Marshal(candidate.request)
		req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		s.mux.ServeHTTP(rr, req)
		resp := rr.Result()

		s.Equalf(candidate.expectedStatusCode, resp.StatusCode, candidate.description)
		if candidate.expectedErr == "" {
			var retrievedUser repository.User
			s.NoErrorf(s.db.First(&retrievedUser, "id = ?", candidate.request.ID).Error, candidate.description)
			s.Equalf(candidate.request, retrievedUser, candidate.description)
		} else {
			bodyBytes, err := io.ReadAll(resp.Body)
			s.NoErrorf(err, candidate.description)
			s.Containsf(string(bodyBytes), candidate.expectedErr, candidate.description)
		}
	}
}

func (s *ServiceTestSuite) TestGetUser() {
	candidates := []struct {
		description        string
		request            string
		expected           repository.User
		expectedErr        string
		expectedStatusCode int
	}{
		// Candidate 1: valid input
		{
			description: "valid input",
			request:     "123e4567-e89b-12d3-a456-426614174000",
			expected: repository.User{
				ID:          "123e4567-e89b-12d3-a456-426614174000",
				Name:        "Alice Johnson",
				Email:       "alice.johnson@example.com",
				DateOfBirth: "1990-04-01T00:00:00Z",
			},
			expectedErr:        "",
			expectedStatusCode: http.StatusOK,
		},
		// Candidate 2: user not found
		{
			description:        "user not found",
			request:            "uuid-not-exists",
			expectedErr:        "404 page not found",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, candidate := range candidates {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", candidate.request), nil)
		rr := httptest.NewRecorder()
		s.mux.ServeHTTP(rr, req)
		resp := rr.Result()

		bodyBytes, err := io.ReadAll(resp.Body)
		s.NoErrorf(err, candidate.description)

		s.Equalf(candidate.expectedStatusCode, resp.StatusCode, candidate.description)
		if candidate.expectedErr == "" {
			var retrievedUser repository.User
			err = json.Unmarshal(bodyBytes, &retrievedUser)
			s.NoErrorf(err, candidate.description)
			s.Equalf(candidate.expected, retrievedUser, candidate.description)
		} else {
			s.Containsf(string(bodyBytes), candidate.expectedErr, candidate.description)
		}
	}
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
