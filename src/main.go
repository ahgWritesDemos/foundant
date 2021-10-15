package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Image struct {
	Filename string
	Body     []byte
}

func (i *Image) save() error {
	filename := "uploads/" + i.Filename // TODO proper path lib?
	return ioutil.WriteFile(filename, i.Body, 0600)
}

func loadImage(filename string) (*Image, error) {
	body, err := ioutil.ReadFile("uploads/" + filename)
	if err != nil {
		return nil, err
	}
	return &Image{Filename: filename, Body: body}, nil
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "welcome to the image server\n")
}

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "I have uploaded your image")
}

func ListImages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "I know about several images")
}

func ShowImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filename := ps.ByName("imageId")
	img, err := loadImage(filename)
	if err != nil {
		fmt.Fprintf(w, "error reading file: %s", err)
		return
	}
	fmt.Fprintf(w, string(img.Body))
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
