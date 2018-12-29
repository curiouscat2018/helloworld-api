package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"

	"golang.org/x/crypto/acme/autocert"
)

const host = "api.curiouscat.one"

type dbEntry struct {
	Greeting     string
	RequestCount int
}

var collection *mongo.Collection

func init() {
	// TODO: check in files and reset connection string
	azureUrl := "mongodb://helloworld-db:CCJF277Zn1Cb0KPD2M579N1JaGsKv1ILOErdqpUDbb4U9XFCkLC6rmF2W1fBmVIz5X7ChzZCswmIHdWahskCwQ==@helloworld-db.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	client, err := mongo.Connect(context.TODO(), azureUrl)
	if err != nil {
		log.Fatalln(err)
	}
	collection = client.Database("main").Collection("main")
	log.Println("Connected to MongoDB")
}

func main() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),
		Cache:      autocert.DirCache("./certs"),
	}

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	http.HandleFunc("/", index)

	log.Println("Start listening helloworld-api")
	if os.Getenv("HELLOWORLD_API_ENV") == "PROD" {
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		log.Fatalln(server.ListenAndServeTLS("", ""))
	} else {
		log.Fatalln(http.ListenAndServe(":http", nil))
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	entry, err := getDBntry()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to insert dbEntry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func getDBntry() (*dbEntry, error) {
	filter := bson.M{"greeting": "helloworld!"}
	update := bson.M{"$inc": bson.M{"requestcount": 1}}

	result := collection.FindOneAndUpdate(context.TODO(), filter, update)
	entry := dbEntry{}
	if err := result.Decode(&entry); err != nil {
		log.Printf("db record not found or corrupted: %v\n", err)

		entry.Greeting = "helloworld!"
		entry.RequestCount = 1
		result, err := collection.InsertOne(context.TODO(), &entry)

		if err != nil {
			return nil, err
		}
		log.Printf("inserted dbEntry %v\n", result.InsertedID)
	}
	return &entry, nil
}
