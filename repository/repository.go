package repository

import (
	"context"
	"database/sql"
	"github.com/stone1549/auth-service/common"
)

// UserRepository represents a data source through which users can be managed.
type UserRepository interface {
	// NewUser adds a user to the repo.
	NewUser(ctx context.Context, email string, password string) (string, error)
	// Authenticate validates email and password combo with what is stored in the repo. Returns users unique id on
	// success
	Authenticate(ctx context.Context, email string, password string) (string, error)
}

// NewUserRepository constructs a UserRepository from the given configuration.
func NewUserRepository(config common.Configuration) (UserRepository, error) {
	var err error
	var repo UserRepository
	var db *sql.DB
	switch config.GetRepoType() {
	case common.InMemoryRepo:
		repo, err = MakeInMemoryRepository(config)
	case common.PostgreSqlRepo:
		db, err = sql.Open("postgres", config.GetPgUrl())

		if err != nil {
			return nil, err
		}
		repo, err = MakePostgresqlUserRespository(config, db)
	default:
		err = newErrRepository("repository type unimplemented")
	}

	return repo, err
}
