package database

import (
	"testing"

	"github.com/curiouscat2018/helloworld-api/config"
	"github.com/stretchr/testify/assert"
)

func TestMockDB_GetDBEntry(t *testing.T) {
	testDB := NewMockDatabase()
	entry, _ := testDB.GetEntry(nil)
	assert.NotEmpty(t, entry.Greeting)
	assert.NotEmpty(t, entry.RequestCount)

	oldCount := entry.RequestCount
	entry, _ = testDB.GetEntry(nil)
	assert.Equal(t, oldCount+1, entry.RequestCount)
}

func TestAzureDB_GetDBEntry(t *testing.T) {
	if testing.Short() {
		return
	}

	testDB, err := NewAzureDatabase(config.Config.DB_URL)
	assert.Nil(t, err)
	entry, err := testDB.GetEntry(nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, entry.Greeting)
	assert.NotEmpty(t, entry.RequestCount)

	oldCount := entry.RequestCount
	entry, err = testDB.GetEntry(nil)
	assert.Nil(t, err)
	assert.Equal(t, oldCount+1, entry.RequestCount)
}
