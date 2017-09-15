package main

import (
	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//

	"log"
	"net/http"

	sw "github.com/AnthonyWyffels/go-server/go"
)

func main() {
	log.Printf("Server started")

	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))

}
