package db

type DBEntry struct {
	Greeting     string
	RequestCount int
}

type DB interface {
	GetDBEntry() (*DBEntry, error)
}
