package database

import "github.com/gin-gonic/gin"

type Entry struct {
	Greeting     string `json:"greeting"`
	RequestCount int `json:"request_count"`
}

type Database interface {
	GetEntry(c *gin.Context) (*Entry, error)
}
