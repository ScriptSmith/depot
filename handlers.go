package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// Job container for index template
type RootData struct {
	Jobs []string
}

type Handlers struct {
	root string
}

// Send an error message to client and log
func logAndRespond(w http.ResponseWriter, r *http.Request, message string) {
	log.Printf("%s: %s", r.RemoteAddr, message)
	http.Error(w, message, 500)
}

// Serve the index page
func (handlers *Handlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		logAndRespond(w, r, "error serving index")
		return
	}

	files, err := ioutil.ReadDir(handlers.root)
	if err != nil {
		logAndRespond(w, r, "error reading root dir")
		return
	}

	dirs := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}

	data := RootData{Jobs: dirs}

	err = tmpl.Execute(w, data)
	if err == nil {
		log.Printf("%s: accessed index", r.RemoteAddr)
	} else {
		logAndRespond(w, r, "error listing jobs")
	}
}

// Serve the job pages
func (handlers *Handlers) JobsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["job"]
	filePath := path.Join(handlers.root, jobName, vars["filepath"])

	_, err := uuid.Parse(jobName)
	if err != nil {
		http.Error(w, "Invalid job name", 400)
		return
	}
	if !path.IsAbs(filePath) {
		http.Error(w, "Invalid file path", 400)
	}

	switch r.Method {
	case "GET":
		log.Printf("%s: serving %s", r.RemoteAddr, filePath)
		http.ServeFile(w, r, filePath)

	case "PUT":
		fileDir, _ := filepath.Split(filePath)
		err = os.MkdirAll(fileDir, 0777)
		if err != nil {
			logAndRespond(w, r, "error creating file directory path")
			return
		}

		tmp, err := ioutil.TempFile("", "put_")
		if err != nil {
			logAndRespond(w, r, "error creating a new file")
			return
		}

		byteCount, err := io.Copy(tmp, r.Body)
		if err != nil {
			logAndRespond(w, r, "error writing file")
			return
		} else if byteCount != r.ContentLength {
			errMsg := fmt.Sprintf("Wrote %d bytes of %d", byteCount, r.ContentLength)
			logAndRespond(w, r, errMsg)
			return
		}

		err = os.Rename(tmp.Name(), filePath)
		if err != nil {
			logAndRespond(w, r, "error adding file at path")
			return
		}

		err = os.Chmod(filePath, 0666)
		if err != nil {
			logAndRespond(w, r, "error updating file permissions")
			return
		} else {
			log.Printf("%s: uploaded %s (%d)", r.RemoteAddr, filePath, byteCount)
		}

	default:
		http.Error(w, "Invalid command", 400)
	}
}

// Serve zip files
func (handlers *Handlers) ZipHandler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	jobName := keys.Get("id")
	if jobName == "" {
		http.Error(w, "No job name provided", 400)
		return
	}

	err := ZipDir(jobName, handlers.root, w, r)
	if err == nil {
		log.Printf("%s: zipped %s", r.RemoteAddr, jobName)
	} else {
		log.Printf("%s: failed zipping %s - %s", r.RemoteAddr, jobName, err)
	}
}
