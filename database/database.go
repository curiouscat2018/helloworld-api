package database

type Entry struct {
	Greeting     string
	RequestCount int
}

type Database interface {
	GetEntry() (*Entry, error)
}
