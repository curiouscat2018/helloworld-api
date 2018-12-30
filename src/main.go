package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/cache"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

const host = "api.curiouscat.one"
const testSecretUrl = "https://helloworld-vault.vault.azure.net/secrets/testsecret/e3246fd47fa74a638e99e4c4afe97006"
const cacheGCSec = 3600
const cachePersistTime = time.Hour

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
	azureUrl := "mongodb://helloworld-db:mMIXPtgqLRa8FWhIzmbuKWTNvSyL2kmdbdewIton3iFp9lqimEhofbMTlQNcNNiSdtmZBfiVpGau5OVLHqPLNg==@helloworld-db.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	client, err := mongo.Connect(context.TODO(), azureUrl)
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

	log.Println("start listening helloworld-api")
	if config.isMockEnv() {
		log.Fatalln(http.ListenAndServe(":http", nil))
	} else {
		go log.Fatalln(http.ListenAndServe(":http", certManager.HTTPHandler(nil)))
		log.Fatalln(server.ListenAndServeTLS("", ""))
	}
}

func index(w http.ResponseWriter, _ *http.Request) {
	entry, err := getDBntry()
	if err != nil {
		log.Println(err)
		reportInternalServerError(w, "failed to get dbEntry", err)
		return
	}

	functor := getSecret
	if config.isMockEnv() {
		functor = getSecretLocal
	}

	secret, err := getSecretFromCache(testSecretUrl, functor)
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

func getSecretFromCache(url string, functor func(string) (string, error)) (string, error) {
	if inmemoryCache.IsExist(url) {
		res, ok := inmemoryCache.Get(url).(string)
		if !ok {
			return "", fmt.Errorf("not able to cast value to string")
		}
		return res, nil
	}

	secert, err := functor(url)
	if err != nil {
		return "", err
	}

	if err := inmemoryCache.Put(url, secert, cachePersistTime); err != nil {
		return "", nil
	}

	return secert, nil
}

func getSecret(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid http status: %v", response.StatusCode)
	}

	responseJSON := struct {
		Value string
	}{}

	if err := json.NewDecoder(response.Body).Decode(&responseJSON); err != nil {
		return "", err
	}

	return responseJSON.Value, nil
}

func reportInternalServerError(w http.ResponseWriter, msg string, err error) {
	log.Printf("%v: %v", msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}
