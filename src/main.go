package main

import (
	"encoding/json"
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
	"image/bmp":  ".bmp",
	"image/gif":  ".gif",
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"video/webm": ".webm",
	"video/mp4":  ".mp4",
}

var allKnownImages = map[string]*Image{}

type Image struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Filename    string    `json:"filename"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// cheap trick to keep the URL tidy
	body, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		fmt.Fprintf(w, "Failed to read static file")
		return
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
	}
	persistFile(file, filename)
	addToIndex(newImg)

	w.Header().Set("Location", "../images/"+uuid.String())
	w.WriteHeader(http.StatusSeeOther)
}

func ListImages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rsp, err := json.Marshal(allKnownImages)
	if err != nil {
		fmt.Fprintf(w, "Failed json marshaling")
		return
	}

	w.Write(rsp)
}

func GetImageMetadata(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	imageId := ps.ByName("imageId")
	img, err := loadImage(imageId)
	if err != nil {
		fmt.Fprintf(w, "error reading file: %s", err)
		return
	}
	jsonResponse(w, img)
}

func GetImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	imageId := ps.ByName("imageId")
	img, err := loadImage(imageId)
	if err != nil {
		fmt.Fprintf(w, "error finding image %s", imageId)
		return
	}

	body, err := ioutil.ReadFile("uploads/" + img.Filename)
	if err != nil {
		fmt.Fprintf(w, "error reading image: %s", err)
		return
	}

	w.Write(body)
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
	return &Image{
			Filename:    img.Filename,
			Id:          img.Id,
			Title:       img.Title,
			Description: img.Description,
		},
		nil
}

func jsonResponse(w http.ResponseWriter, v interface{}) {
	rsp, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(w, "Failed json marshaling")
		return
	}

	w.Write(rsp)
}

/*--------------------------------------------*/

func main() {

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/upload/", Upload)
	router.GET("/images/", ListImages)
	router.GET("/images/:imageId", GetImage)
	router.GET("/images/:imageId/metadata", GetImageMetadata)
	router.ServeFiles("/static/*filepath", http.Dir("./static/"))

	log.Println("Starting Server on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
