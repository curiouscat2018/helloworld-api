package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// certManager := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist("wwww.curiouscat.one"), //Your domain here
	// 	Cache:      autocert.DirCache("certs"),                    //Folder for storing certificates
	// }

	// server := &http.Server{
	// 	Addr: ":https",
	// 	TLSConfig: &tls.Config{
	// 		GetCertificate: certManager.GetCertificate,
	// 	},
	// }

	http.HandleFunc("/", index)
	//log.Fatal(server.ListenAndServeTLS("", "")) //Key and cert are coming from Let's Encrypt
	log.Fatal(http.Serve(autocert.NewListener("www.curiouscat.one"), nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p>Hello world!</p>")
}
