package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

type apiResponse struct {
	Data string `json:"data"`
}

func main() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("www.curiouscat.one"),
		Cache:      autocert.DirCache("./certs"),
	}

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	http.HandleFunc("/", index)

	log.Println("start helloworld-api")
	if os.Getenv("HELLOWORLD_API_ENV") == "PROD" {
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		log.Fatalln(server.ListenAndServeTLS("", ""))
	} else {
		log.Fatalln(http.ListenAndServe(":http", nil))
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	res := apiResponse{Data: "Hello world!"}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		log.Println("Failed to encode JSON")
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
