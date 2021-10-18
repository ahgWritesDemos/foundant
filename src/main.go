package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type Image struct {
	Title       string
	Description string
	Filename    string
	Body        []byte
}

func (i *Image) save(fileHeader *multipart.FileHeader) error {
	// Given a fleshed-out image, save it to some filename with the right prefix, and note the filename
	source, err := fileHeader.Open()
	if err != nil {
		panic("failed to open parsed form file" + err.Error())
	}

	destination, err := os.Create("uploads/" + i.Filename)
	if err != nil {
		panic("failed to create destination" + err.Error())
	}
	len, err := io.Copy(destination, source)
	if err != nil {
		panic("failed on copy" + err.Error())
	}
	log.Print("copied %d bytes", len)

	return nil
}

func loadImage(filename string) (*Image, error) {
	body, err := ioutil.ReadFile("uploads/" + filename)
	if err != nil {
		return nil, err
	}
	return &Image{Filename: filename, Body: body}, nil
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// cheap trick to keep the URL tidy
	body, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		fmt.Fprintf(w, "Failed to read static file")
	}
	w.Write(body)
}

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(32 << 20)
	file := r.MultipartForm.File["upload"][0]
	filename := generateFilename(file)

	newImg := &Image{
		Filename:    filename,
		Title:       r.Form["title"][0],
		Description: r.Form["description"][0],
		Body:        nil,
	}
	newImg.save(file)

	log.Print("uploaded %v", newImg)

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
	w.Write(img.Body)
}

func main() {

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/upload/", Upload)
	router.GET("/images/", ListImages)
	router.GET("/images/:imageId", ShowImage)
	router.ServeFiles("/static/*filepath", http.Dir("./static/"))

	log.Println("Starting Server on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Declare a map from MIME type to a canonical extension
//   so that I can use arbitrary filenames
var canonicalExtensions = map[string]string{
	// <form accept=".png,.jpg,.jpeg,.gif,.mp4,.webm,.bmp" ... />
	"image/gif": ".gif",
	"image/jpeg": ".jpg",
}

func generateFilename(fh *multipart.FileHeader) string {
	ext := canonicalExtensions["image/jpeg"]
	root := uuid.New()
	return root.String() + ext
}
