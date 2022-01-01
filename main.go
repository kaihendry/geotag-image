package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"strconv"

	"github.com/apex/gateway/v2"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/tajtiattila/metadata/exif"

	jsonhandler "github.com/apex/log/handlers/json"
)

//go:embed public
var public embed.FS

func (s *server) upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseMultipartForm(0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(r.Form["lat"]) == 0 {
			http.Error(w, "missing latitude", http.StatusBadRequest)
			return
		}

		if len(r.Form["lng"]) == 0 {
			http.Error(w, "missing longitude", http.StatusBadRequest)
			return
		}

		lat, err := strconv.ParseFloat(r.Form["lat"][0], 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		lng, err := strconv.ParseFloat(r.Form["lng"][0], 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.WithFields(log.Fields{
			"lat": lat,
			"lng": lng,
		}).Info("geolocation")

		file, _, err := r.FormFile("jpeg")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		x, err := exif.Decode(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		x.SetLatLong(lat, lng)

		_, err = file.Seek(0, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "image/jpeg")

		if err := exif.Copy(w, file, x); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

type server struct {
	router *http.ServeMux
}

func newServer(local bool) *server {

	s := &server{router: &http.ServeMux{}}

	if local {
		log.SetHandler(text.Default)
	} else {
		log.SetHandler(jsonhandler.Default)
	}

	directory, err := fs.Sub(public, "public")
	if err != nil {
		log.WithError(err).Fatal("unable to load public static files")
	}
	fileServer := http.FileServer(http.FS(directory))
	s.router.Handle("/public/", http.StripPrefix("/public", fileServer))
	s.router.Handle("/upload", s.upload())
	s.router.Handle("/", s.index())

	return s
}

func main() {

	_, awsDetected := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	log.WithField("awsDetected", awsDetected).Info("starting up")

	s := newServer(!awsDetected)

	var err error

	if awsDetected {
		log.Info("starting cloud server")
		err = gateway.ListenAndServe("", s.router)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), s.router)
	}
	log.WithError(err).Fatal("error listening")

}

func (s *server) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("").ParseFS(public, "public/index.html"))
		w.Header().Add("Content-Type", "text/html")
		err := t.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithError(err).Fatal("Failed to execute templates")
		}
	}
}
