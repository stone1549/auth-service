package repository

import (
	"context"
	"encoding/json"
	"github.com/stone1549/auth-service/common"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"time"
)

type storedUser struct {
	common.User
	Id         string    `json:"id"`
	SaltedHash string    `json:"saltedHash"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
type inMemoryUserRepository struct {
	usersByEmail map[string]*storedUser
}

// NewUser adds a user to the repo.
func (imr *inMemoryUserRepository) NewUser(ctx context.Context, email string, password string) (string, error) {
	if email == "" {
		return "", newErrRepository("email is required")
	} else if password == "" {
		return "", newErrRepository("password is required")
	}

	id := uuid.NewV4().String()

	_, ok := imr.usersByEmail[email]
	if ok {
		return "", newErrRepository("user already exists")
	}

	saltedHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", newErrRepository("unable to generate password")
	}

	createdAt := time.Now()
	updatedAt := createdAt

	imr.usersByEmail[email] = &storedUser{common.User{Email: email}, id, string(saltedHash),
		createdAt, updatedAt}

	return id, nil
}

// Authenticate compares the given email and password combination against the salted hash in the repo.
func (imr *inMemoryUserRepository) Authenticate(ctx context.Context, email string, password string) (string, error) {
	if email == "" {
		return "", newErrRepository("email is required")
	} else if password == "" {
		return "", newErrRepository("password is required")
	}

	user, ok := imr.usersByEmail[email]
	if !ok {
		return "", newErrRepository("user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.SaltedHash), []byte(password)) != nil {
		return "", newErrRepository("invalid username/password combo")
	}

	return user.Id, nil
}

// MakeInMemoryRepository constructs an in memory backed UserRepository from the given configuration.
func MakeInMemoryRepository(config common.Configuration) (UserRepository, error) {
	var err error

	usersByEmail, err := loadInitInMemoryDataset(config.GetInitDataSet())

	return &inMemoryUserRepository{usersByEmail}, err
}

func loadInitInMemoryDataset(dataset string) (map[string]*storedUser, error) {
	if dataset == "" {
		return make(map[string]*storedUser), nil
	}

	var err error
	storedUsers := make([]storedUser, 0)

	if err != nil {
		return nil, err
	}

	jsonBytes, err := ioutil.ReadFile(dataset)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &storedUsers)

	if err != nil {
		return nil, err
	}

	usersByEmail := make(map[string]*storedUser)

	for index, storedUser := range storedUsers {
		usersByEmail[storedUser.Email] = &storedUsers[index]
	}

	return usersByEmail, err
}
