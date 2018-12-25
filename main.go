package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Go server started.")
	http.HandleFunc("/", index)
	go func() {
		log.Fatalln(http.ListenAndServe(":80", nil))
	}()

	log.Fatalln(http.ListenAndServeTLS(":443", "./cert.pem", "./key.pem", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p>Hello world!</p>")
}
