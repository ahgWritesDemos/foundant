package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "welcome\n")
}

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "ack")
}
func ListImages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "I know about several images")
}
func ShowImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "here is image %s\n", ps.ByName("imageId"))
}

func main() {

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/images/", Upload)
	router.GET("/images/", ListImages)
	router.GET("/images/:imageId", ShowImage)

	log.Println("Starting Server on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
