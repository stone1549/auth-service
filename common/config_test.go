package common_test

import (
	"github.com/stone1549/auth-service/common"
	"os"
	"testing"
)

const (
	lifeCycleKey       string = "AUTH_SERVICE_ENVIRONMENT"
	repoTypeKey        string = "AUTH_SERVICE_REPO_TYPE"
	timeoutSecondsKey  string = "AUTH_SERVICE_TIMEOUT"
	portKey            string = "AUTH_SERVICE_PORT"
	pgUrlKey           string = "AUTH_SERVICE_PG_URL"
	pgInitDatasetKey   string = "AUTH_SERVICE_INIT_DATASET"
	tokenSecretKeyKey  string = "AUTH_SERVICE_TOKEN_SECRET"
	tokenPrivateKeyKey string = "AUTH_SERVICE_TOKEN_PRIV"
	tokenPublicKeyKey  string = "AUTH_SERVICE_TOKEN_PUB"
)

func clearEnv() {
	os.Setenv(lifeCycleKey, "")
	os.Setenv(repoTypeKey, "")
	os.Setenv(timeoutSecondsKey, "")
	os.Setenv(portKey, "")
	os.Setenv(pgUrlKey, "")
	os.Setenv(pgInitDatasetKey, "")
	os.Setenv(tokenSecretKeyKey, "")
	os.Setenv(tokenPrivateKeyKey, "../data/sample.key")
	os.Setenv(tokenPublicKeyKey, "../data/sample.pub")
}

func setEnv(lifeCycle, repoType, timeoutSeconds, port, pgUrl, pgInitDataset, tokenSecretKey, tokenPrivateKey,
	tokenPublicKey string) {
	os.Setenv(lifeCycleKey, lifeCycle)
	os.Setenv(repoTypeKey, repoType)
	os.Setenv(timeoutSecondsKey, timeoutSeconds)
	os.Setenv(portKey, port)
	os.Setenv(pgUrlKey, pgUrl)
	os.Setenv(pgInitDatasetKey, pgInitDataset)
	os.Setenv(tokenSecretKeyKey, tokenSecretKey)
	os.Setenv(tokenPrivateKeyKey, tokenPrivateKey)
	os.Setenv(tokenPublicKeyKey, tokenPublicKey)
}

// TestGetConfiguration_Defaults ensures that a default configuration is returned if no configuration is provided in
// the environment.
func TestGetConfiguration_Defaults(t *testing.T) {
	clearEnv()
	_, err := common.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_Defaults ensures that a default configuration is returned if no configuration is provided in
// the environment.
func TestGetConfiguration_ImSuccess(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "", "",
		"", "", "")
	_, err := common.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_ImSuccessSmallDataset ensures that a configuration is returned when specifying an in memory
// repo with an initial dataset.
func TestGetConfiguration_ImSuccessSmallDataset(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "",
		"../data/small_set.json", "SECRET", "", "")
	_, err := common.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_ImSuccessNoneDataset ensures that a configuration is returned when specifying an in memory
// repo without an initial dataset.
func TestGetConfiguration_ImSuccessNoneDataset(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "", "",
		"SECRET", "", "")
	_, err := common.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_FailRepo ensures that an error is returned when specifying an invalid repo type.
func TestGetConfiguration_FailRepo(t *testing.T) {
	setEnv("PROD", "", "60", "3333", "", "",
		"SECRET", "", "")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_FailTimeout ensures that an error is returned when specifying an invalid timeout.
func TestGetConfiguration_FailTimeout(t *testing.T) {
	setEnv("PROD", "IN_MEMORY", "", "3333", "", "",
		"SECRET", "", "")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_FailPort ensures that an error is returned when specifying an invalid port.
func TestGetConfiguration_FailPort(t *testing.T) {
	setEnv("PROD", "IN_MEMORY", "60", "", "", "",
		"SECRET", "", "")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_PgSuccess ensures that a configuration is returned when specifying a PostgreSQL repo type.
func TestGetConfiguration_PgSuccess(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333",
		"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "",
		"SECRET", "", "")
	_, err := common.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_PgFailPgUrl ensures that an error is returned when specifying a PostgreSQL repo type without a
// connection url.
func TestGetConfiguration_PgFailPgUrl(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333", "", "",
		"SECRET", "", "")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_ImFailNoJwtKey ensures that an error is returned when no JWT key is provided.
func TestGetConfiguration_ImFailNoJwtKey(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333", "", "",
		"", "", "")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_ImFailRsaJustPrivateKey ensures that an error is returned when no public JWT key is provided.
func TestGetConfiguration_ImFailRsaJustPrivateKey(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333", "", "",
		"", "asdfdsaf", "")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_ImFailRsaJustPublicKey ensures that an error is returned when no public JWT key is provided.
func TestGetConfiguration_ImFailRsaJustPublicKey(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333", "", "",
		"", "", "asdfdsaf")
	_, err := common.GetConfiguration()
	notOk(t, err)
}
