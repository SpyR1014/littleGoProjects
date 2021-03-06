package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Process a multipart form.
func process(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if len(r.MultipartForm.File["uploaded"]) < 1 {
		http.Error(w, "Error: Please submit a file!", 400)
		return
	}
	fileHeader := r.MultipartForm.File["uploaded"][0]
	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintln(w, string(data))

}

// This function returns the form.
func serveIndex(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("index.html")
	w.Write(file)
}

// Example of returning different status codes.
func noSuchFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	fmt.Fprintln(w, "No such service...")
}

// Example of how headers can = redirect.
func headerRedirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "http://google.com")
	w.WriteHeader(301)
}

type post struct {
	User    string
	Threads []string
}

func serveJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := &post{
		User:    "Andrew J",
		Threads: []string{"first", "second", "third"},
	}
	json, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Failed to marshal Post to JSON", 500)
		return
	}
	w.Write(json)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/process", process)
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/none", noSuchFunction)
	http.HandleFunc("/escape", headerRedirect)
	http.HandleFunc("/json", serveJSON)

	fmt.Println("Starting server on [", server.Addr, "]")
	server.ListenAndServe()
}
