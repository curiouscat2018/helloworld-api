package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/curiouscat2018/helloworld-api/config"
	"github.com/curiouscat2018/helloworld-api/database"
)

var myDB database.Database

func main() {
	http.HandleFunc("/", index)
	log.Printf("start listening helloworld-api: isMockEnv: %v", config.Config.IsMockEnv)

	if config.Config.IsMockEnv {
		prepareMockEnv(&myDB)
	} else {
		prepareProdEnv(&myDB)
	}
	log.Fatalln(http.ListenAndServe(":http", nil))
}

func prepareProdEnv(db *database.Database) {
	log.Println("preparing Azure Database")
	tempDB, err := database.NewAzureDatabase(config.Config.DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	*db = tempDB
}

func prepareMockEnv(db *database.Database) {
	log.Println("preparing mock Database")
	tempDB := database.NewMockDatabase()
	*db = tempDB
}

func index(w http.ResponseWriter, _ *http.Request) {
	entry, err := myDB.GetEntry()
	if err != nil {
		log.Println(err)
		reportInternalServerError(w, "failed to get database entry", err)
		return
	}

	response := struct {
		database.Entry
		HostName string
	}{
		Entry:    *entry,
		HostName: config.Config.HostName(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		reportInternalServerError(w, "failed to encode JSON", err)
	}
	log.Print("handle index success")
}

func reportInternalServerError(w http.ResponseWriter, msg string, err error) {
	log.Printf("%v: %v", msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}
