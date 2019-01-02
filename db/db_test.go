package db

import (
	"testing"

	"github.com/curiouscat2018/helloworld-api/config"
	"github.com/curiouscat2018/helloworld-api/vault"
	"github.com/stretchr/testify/assert"
)

func TestMockDB_GetDBEntry(t *testing.T) {
	testDB, _ := NewMockDB()
	entry, _ := testDB.GetDBEntry()
	assert.NotEmpty(t, entry.Greeting)
	assert.NotEmpty(t, entry.RequestCount)

	oldCount := entry.RequestCount
	entry, _ = testDB.GetDBEntry()
	assert.Equal(t, oldCount+1, entry.RequestCount)
}

func TestAzureDB_GetDBEntry(t *testing.T) {
	if testing.Short() {
		return
	}

	v := vault.NewAzureVault()
	dbURL, err := v.GetSecret(config.DBURLVaultIdentifier)
	assert.Nil(t, err)

	testDB, err := NewAzureDB(dbURL)
	assert.Nil(t, err)
	entry, err := testDB.GetDBEntry()
	assert.Nil(t, err)
	assert.NotEmpty(t, entry.Greeting)
	assert.NotEmpty(t, entry.RequestCount)

	oldCount := entry.RequestCount
	entry, err = testDB.GetDBEntry()
	assert.Nil(t, err)
	assert.Equal(t, oldCount+1, entry.RequestCount)
}
