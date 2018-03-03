package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/apex/log"
	"github.com/gorilla/pat"
	"github.com/tajtiattila/metadata/exif"
)

func upload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(0)

	r.ParseForm()

	fmt.Println("Lat:", r.Form["lat"])
	fmt.Println("Lng:", r.Form["lng"])
	lat, err := strconv.ParseFloat(r.Form["lat"][0], 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	lng, err := strconv.ParseFloat(r.Form["lng"][0], 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	file, _, err := r.FormFile("jpeg")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	filetype := http.DetectContentType(buff)
	fmt.Println(filetype)

	if filetype != "image/jpeg" {
		http.Error(w, "Upload not a JPEG", 400)
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	x, err := exif.Decode(file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	x.SetLatLong(lat, lng)

	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")

	if err := exif.Copy(w, file, x); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func getStatic(w http.ResponseWriter, r *http.Request) {
	// log.Infof("Requested: %s", r.URL.Path[1:])
	http.ServeFile(w, r, "dist/"+r.URL.Path[1:])
}

func main() {

	addr := ":" + os.Getenv("PORT")
	app := pat.New()

	app.Post("/", upload)
	app.Get("/js/", http.HandlerFunc(getStatic))
	app.Get("/css/", http.HandlerFunc(getStatic))
	app.Get("/favicon.ico", http.HandlerFunc(getStatic))
	app.Get("/", index)

	if err := http.ListenAndServe(addr, app); err != nil {
		log.WithError(err).Fatal("error listening")
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("").ParseGlob("dist/*.html"))
	t.ExecuteTemplate(w, "index.html", nil)
}
