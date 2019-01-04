package database

type mockDB struct {
	count int
}

func NewMockDatabase() Database {
	db := &mockDB{}
	db.count = 0
	return db
}

func (db *mockDB) GetEntry() (*Entry, error) {
	db.count++
	entry := Entry{
		Greeting:     "helloworld! from mock Database",
		RequestCount: db.count,
	}
	return &entry, nil
}
