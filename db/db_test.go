package db

import (
	"../config"
	"flag"
	"github.com/stretchr/testify/assert"
	"testing"
)

var integrationTest = flag.Bool("int", false, "run integration test")

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
	if !*integrationTest {
		return
	}

	testDB, err := NewAzureDB(config.DBURL)
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
