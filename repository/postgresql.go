package repository

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stone1549/auth-service/common"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	insertLogin       = "INSERT INTO login (id, email, salted_hash) VALUES ($1, $2, $3)"
	authenticate      = "SELECT salted_hash, id FROM login WHERE email=$1"
	insertStoredLogin = "INSERT INTO login (id, email, salted_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
)

type postgresqlUserRepository struct {
	db *sql.DB
}

// NewUser adds a user to the repo.
func (impr *postgresqlUserRepository) NewUser(ctx context.Context, email string, password string) (string, error) {
	if email == "" {
		return "", newErrRepository("email is required")
	}

	id := uuid.NewV4().String()

	saltedHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", newErrRepository("unable to generate password")
	}

	_, err = impr.db.ExecContext(ctx, insertLogin, id, email, saltedHash)

	return id, err
}

// Authenticate compares a given email and password combination against the salted hash in the repo.
func (impr *postgresqlUserRepository) Authenticate(ctx context.Context, email string, password string) (string, error) {
	if email == "" {
		return "", newErrRepository("email is required")
	}

	if password == "" {
		return "", newErrRepository("password is required")
	}

	row := impr.db.QueryRowContext(ctx, authenticate, email)
	var saltedHash string
	var id string

	err := row.Scan(&saltedHash, &id)

	if err == sql.ErrNoRows {
		return "", newErrRepository("user not found")
	} else {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(saltedHash), []byte(password))

	if err != nil {
		return "", errors.New("invalid username/password combo")
	}

	return id, nil
}

func loadInitPostgresqlData(db *sql.DB, dataset string) error {
	users, err := loadInitInMemoryDataset(dataset)

	if err != nil {
		return err
	}

	txn, err := db.Begin()

	if err != nil {
		return err
	}

	for id, user := range users {
		_, err = txn.Exec(insertStoredLogin, id, user.Email, user.SaltedHash, user.CreatedAt, user.UpdatedAt)

		if err != nil {
			return err
		}
	}

	return txn.Commit()
}

// MakePostgresqlUserRespository constructs a PostgreSQL backed UserRepository from the given params.
func MakePostgresqlUserRespository(config common.Configuration, db *sql.DB) (UserRepository, error) {
	var err error
	if config.GetInitDataSet() != "" {
		err = loadInitPostgresqlData(db, config.GetInitDataSet())
	}

	if err != nil {
		return nil, err
	}

	return &postgresqlUserRepository{db}, nil
}
