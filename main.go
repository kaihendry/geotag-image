package main

import (
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/apex/log"
	"github.com/gorilla/pat"
	"github.com/tajtiattila/metadata/exif"
)

func upload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// log.Infof("%+v", r.Form)
	// log.Infof("%+v", r.PostForm)

	if len(r.Form["lat"]) == 0 {
		http.Error(w, "missing latitude", 400)
		return
	}

	if len(r.Form["lng"]) == 0 {
		http.Error(w, "missing longitude", 400)
		return
	}

	log.Infof("Lat: %v Lng: %v", r.Form["lat"], r.Form["lng"])

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
	log.Infof("Upload filetype: %s", filetype)

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
