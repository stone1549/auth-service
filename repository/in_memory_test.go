package repository_test

import (
	"github.com/stone1549/auth-service/repository"
	"testing"
)

func makeNewImRepo(t *testing.T) repository.UserRepository {
	repo, err := repository.MakeInMemoryRepository(inMemorySmall)

	ok(t, err)
	return repo
}
