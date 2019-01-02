package main

import (
	"./cache"
	. "./config"
	"./db"
	"./vault"
	"crypto/tls"
	"encoding/json"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
)

var myCache *cache.Cache
var myVault vault.Vault
var myDB db.DB

func main() {
	http.HandleFunc("/", index)
	log.Printf("start listening helloworld-api: isMockEnv: %v", Config.IsMockEnv())

	if Config.IsMockEnv() {
		myCache, myVault, myDB = prepareMockEnv()
		log.Fatalln(http.ListenAndServe(":http", nil))
	} else {
		myCache, myVault, myDB = prepareProdEnv()
		certManager, server := prepareTLS()
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		log.Fatalln(server.ListenAndServeTLS("", ""))
	}
}

func prepareProdEnv() (c *cache.Cache, v vault.Vault, d db.DB) {
	d, err := db.NewAzureDB(DBURL)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("prepare Azure DB")

	c, err = cache.NewCache(cache.CacheGCSec, cache.CachePersistTime)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("prepare in-memory cache")

	v = vault.NewAzureVault()
	log.Println("prepare Azure vault")
	return
}

func prepareMockEnv() (c *cache.Cache, v vault.Vault, d db.DB) {
	d, err := db.NewMockDB()
	log.Println("prepare mock DB")

	c, err = cache.NewCache(cache.CacheGCSec, cache.CachePersistTime)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("prepare in-memory cache")

	v = vault.NewMockVault()
	log.Println("prepare mock vault")
	return
}

func prepareTLS() (autocert.Manager, *http.Server) {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(Host),
		Cache:      autocert.DirCache("./certs"),
	}

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	return certManager, server
}

func index(w http.ResponseWriter, _ *http.Request) {
	entry, err := myDB.GetDBEntry()
	if err != nil {
		log.Println(err)
		reportInternalServerError(w, "failed to get dbEntry", err)
		return
	}

	secret, _, err := myCache.GetRESTDataFromCache(DemoSecretUrl, myVault.GetSecret)
	if err != nil {
		reportInternalServerError(w, "failed to get demosecret", err)
		return
	}

	response := struct {
		db.DBEntry
		DemoSecret string
	}{
		*entry,
		secret,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		reportInternalServerError(w, "failed to encode JSON", err)
	}
}

func reportInternalServerError(w http.ResponseWriter, msg string, err error) {
	log.Printf("%v: %v", msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}
