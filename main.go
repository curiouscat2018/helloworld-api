package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"

	"github.com/curiouscat2018/helloworld-api/cache"
	"github.com/curiouscat2018/helloworld-api/config"
	"github.com/curiouscat2018/helloworld-api/db"
	"github.com/curiouscat2018/helloworld-api/vault"
	"golang.org/x/crypto/acme/autocert"
)

var myCache *cache.Cache
var myVault vault.Vault
var myDB db.DB

func main() {
	http.HandleFunc("/", index)
	log.Printf("start listening helloworld-api: isMockEnv: %v", config.Config.IsMockEnv())

	if config.Config.IsMockEnv() {
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
	log.Println("preparing in-memory cache")
	c, err := cache.NewCache(cache.CacheGCSec, cache.CachePersistTime)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("preparing Azure vault")
	v = vault.NewAzureVault()

	log.Println("preparing Azure DB")
	dbURL, err := v.GetSecret(config.DBURLVaultIdentifier)
	if err != nil {
		log.Fatal(err)
	}

	d, err = db.NewAzureDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func prepareMockEnv() (c *cache.Cache, v vault.Vault, d db.DB) {
	log.Println("preparing in-memory cache")
	c, err := cache.NewCache(cache.CacheGCSec, cache.CachePersistTime)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("preparing mock vault")
	v = vault.NewMockVault()

	log.Println("preparing mock DB")
	d, err = db.NewMockDB()
	return
}

func prepareTLS() (autocert.Manager, *http.Server) {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Host),
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

	secret, _, err := myCache.GetRESTDataFromCache(config.DemosecretVaultIdentifier, myVault.GetSecret)
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
