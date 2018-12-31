package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/astaxie/beego/cache"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"strconv"
)

const host = "api.curiouscat.one"
const testSecretUrl = "https://helloworld-vault.vault.azure.net/secrets/testsecret/e3246fd47fa74a638e99e4c4afe97006?api-version=7.0"


type dbEntry struct {
	Greeting     string
	RequestCount int
}

type apiResponse struct {
	dbEntry
	TestSecret string
}

var collection *mongo.Collection
var inmemoryCache cache.Cache

func init() {
	// TODO: check in files and reset connection string
	url := "mongodb://helloworld-db:mMIXPtgqLRa8FWhIzmbuKWTNvSyL2kmdbdewIton3iFp9lqimEhofbMTlQNcNNiSdtmZBfiVpGau5OVLHqPLNg==@helloworld-db.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	client, err := mongo.Connect(context.TODO(), url)
	if err != nil {
		log.Fatalln(err)
	}
	collection = client.Database("main").Collection("main")
	log.Println("connected to MongoDB")

	inmemoryCache, err = cache.NewCache("memory", `{"interval":`+strconv.Itoa(cacheGCSec)+`}`)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("created new in-memory cache")
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

	log.Printf("home: %v", config.getHomeDir())
	if config.isMockEnv() {
		log.Println("start listening mock helloworld-api")
		log.Fatalln(http.ListenAndServe(":http", nil))
	} else {
		log.Println("start listening prod helloworld-api")
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		log.Fatalln(server.ListenAndServeTLS("", ""))
	}
}

func index(w http.ResponseWriter, _ *http.Request) {
	entry, err := getDBEntry()
	if err != nil {
		log.Println(err)
		reportInternalServerError(w, "failed to get dbEntry", err)
		return
	}

	functor := getSecret
	if config.isMockEnv() {
		functor = getSecretLocal
	}

	secret, err := getRESTDataFromCache(testSecretUrl, functor)
	if err != nil {
		reportInternalServerError(w, "failed to get testSecret", err)
		return
	}

	response := apiResponse{
		dbEntry:    *entry,
		TestSecret: secret,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		reportInternalServerError(w, "failed to encode JSON", err)
	}
}

func getDBEntry() (*dbEntry, error) {
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

func reportInternalServerError(w http.ResponseWriter, msg string, err error) {
	log.Printf("%v: %v", msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}
