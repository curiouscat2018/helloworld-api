package database

type Entry struct {
	Greeting     string `json:"greeting"`
	RequestCount int `json:"request_count"`
}

type Database interface {
	GetEntry() (*Entry, error)
}
