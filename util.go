package main

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func ValidateMedia(w http.ResponseWriter, handler *multipart.FileHeader) bool {
	if handler.Header.Get("Content-Type") != "image/jpeg" && handler.Header.Get("Content-Type") != "image/jpg" {
		http.Error(w, "wrong content type:"+handler.Header.Get("Content-Type")+", expected jpeg or jpg.", http.StatusUnsupportedMediaType)
		return true
	}
	if handler.Size > 10485760 {
		http.Error(w, "size too big"+strconv.FormatInt(handler.Size, 10)+" byte, max size is 10 mb", http.StatusBadRequest)
		return true
	}
	return false
}


func ResizePic(w http.ResponseWriter, err error, proportion int, fileName string) string {

	fileNameSmaller := fileName + "-" + strconv.Itoa(proportion) + "percent.jpg"

	file, err := os.Open("./tmp/" + fileName + ".jpg")
	img, err := jpeg.Decode(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ratio := (float64(100) - float64(proportion)) / float64(100)
	m := resize.Resize(uint(float64(img.Bounds().Max.Y)*ratio), 0, img, resize.Lanczos3)

	out, err := os.Create("./tmp/" + fileNameSmaller)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jpeg.Encode(out, m, nil)
	out.Close()
	file.Close()

	return fileNameSmaller
}