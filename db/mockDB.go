package db

type mockDB struct {
	count int
}

func NewMockDB() (DB, error) {
	db := &mockDB{}
	db.count = 0
	return db, nil
}

func (db *mockDB) GetDBEntry() (*DBEntry, error) {
	db.count++
	entry := DBEntry{
		Greeting:     "helloworld! from mock DB",
		RequestCount: db.count,
	}
	return &entry, nil
}
