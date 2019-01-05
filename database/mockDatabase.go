package database

import "github.com/gin-gonic/gin"

type mockDB struct {
	count int
}

func NewMockDatabase() Database {
	db := &mockDB{}
	db.count = 0
	return db
}

func (db *mockDB) GetEntry(c *gin.Context) (*Entry, error) {
	db.count++
	entry := Entry{
		Greeting:     "helloworld! from mock Database",
		RequestCount: db.count,
	}
	return &entry, nil
}
