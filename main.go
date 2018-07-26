package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg := "Hello World, Value from config => " + os.Getenv("CONFIG_VALUE")
	io.WriteString(w, msg)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
