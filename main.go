package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("wwww.curiouscat.one", "curiouscat.one"),
		Cache:      autocert.DirCache("./certs"),
	}

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	http.HandleFunc("/", index)

	if os.Getenv("HELLOWORLD_SERVER_ENV") == "PROD" {
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		http.ListenAndServe(":http", nil)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p>Hello world!</p>")
}
