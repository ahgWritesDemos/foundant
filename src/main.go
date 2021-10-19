package main

import (
	"errors"
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

// Declare a map from MIME type to a canonical extension
//   so that I can use arbitrary filenames
var canonicalExtensions = map[string]string{
	// <form accept=".png,.jpg,.jpeg,.gif,.mp4,.webm,.bmp" ... />
	"image/gif":  ".gif",
	"image/jpeg": ".jpg",
}

var allKnownImages = map[string]*Image{}

type Image struct {
	Id          uuid.UUID
	Title       string
	Description string
	Filename    string
	Body        []byte
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
	uuid, filename := generateFilename(file)

	newImg := &Image{
		Id:          uuid,
		Filename:    filename,
		Title:       r.Form["title"][0],
		Description: r.Form["description"][0],
		Body:        nil,
	}
	persistFile(file, filename)
	addToIndex(newImg)

	w.WriteHeader(http.StatusSeeOther)
	w.Header().Set("Location", "images/"+uuid.String())
}

func ListImages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "I know about %d images", len(allKnownImages))
}

func ShowImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	imageId := ps.ByName("imageId")
	img, err := loadImage(imageId)
	if err != nil {
		fmt.Fprintf(w, "error reading file: %s", err)
		return
	}
	w.Write(img.Body)
}

/*--------------------------------------------*/

func persistFile(fileHeader *multipart.FileHeader, filename string) error {
	source, err := fileHeader.Open()
	if err != nil {
		return err
	}
	destination, err := os.Create("uploads/" + filename)
	if err != nil {
		return err
	}

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return nil
}

func addToIndex(img *Image) {
	allKnownImages[img.Id.String()] = img
}

func generateFilename(fh *multipart.FileHeader) (uuid.UUID, string) {
	ext := canonicalExtensions["image/jpeg"]
	root := uuid.New()
	return root, root.String() + ext
}

func loadImage(id string) (*Image, error) {
	img := allKnownImages[id]
	if img == nil {
		return nil, errors.New("no such image")
	}
	filename := img.Filename

	body, err := ioutil.ReadFile("uploads/" + filename)
	if err != nil {
		return nil, err
	}
	return &Image{Filename: filename, Body: body}, nil
}

/*--------------------------------------------*/

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
