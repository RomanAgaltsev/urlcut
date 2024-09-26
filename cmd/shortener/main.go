package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	srvAddr string = ":8080"
	urlShrt string = "http://localhost:8080/EwHXdJfB"
)

var url []byte

func basicHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/":
		postHandler(w, r)
	case r.Method == http.MethodGet && r.URL.Path != "/":
		getHandler(w, r)
	default:
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	url, _ = io.ReadAll(r.Body)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(urlShrt)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(urlShrt))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", string(url))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicHandler)
	log.Fatal(http.ListenAndServe(srvAddr, mux))
}
