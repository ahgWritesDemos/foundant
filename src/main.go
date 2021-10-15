package main

import (
	"log"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "welcome\n")
}

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func main() {

	router := httprouter.New()
	router.GET("/", Index)

	log.Println("Starting Server on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
