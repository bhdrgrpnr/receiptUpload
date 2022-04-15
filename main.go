package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Compile templates on start of the application
var templates = template.Must(template.ParseFiles("public/upload.html"))

// Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {

	fileName := strconv.Itoa(rand.Int())

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ValidateMedia(w, handler) {
		return
	}

	fileOriginal := "./tmp/" + fileName + ".jpg"
	dst, err := os.Create(fileOriginal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file.Close()
	dst.Close()

	fileSize90 := ResizePic(w, err, 90, fileName)
	fileSize50 := ResizePic(w, err, 50, fileName)

	printOutput(w, fileName, fileSize50, fileSize90)
}

func printOutput(w http.ResponseWriter, fileName string, fileSize50 string, fileSize90 string) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%s", "Successfully Uploaded File:\n <a href=\"  "+fileName+".jpg\">"+fileName+".jpg</a></br> ")
	fmt.Fprintf(w, "%s", "Successfully Uploaded File:\n <a href=\"  "+fileSize50+"  \">"+fileSize50+"</a></br> ")
	fmt.Fprintf(w, "%s", "Successfully Uploaded File:\n <a href=\"  "+fileSize90+"  \">"+fileSize90+"</a></br> ")
}



func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}

func main() {

	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/", http.FileServer(http.Dir("./tmp")))

	http.ListenAndServe(":8080", nil)
}
