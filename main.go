package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg := fmt.Sprintf("Hello World from %s\n\n", os.Getenv("HOSTNAME"))
	msg = msg + "Env vars:\n\n"
	for _, p := range os.Environ() {
		msg = msg + fmt.Sprintf("%s\n", p)
	}
	io.WriteString(w, msg)
	generateLog(r.Host)
}

func generateLog(msg string) {
	lg := fmt.Sprintf("hit from %s, logID=%s", msg, uuid.New().String())
	logStdout := log.New(os.Stdout, "stdout", log.LstdFlags)
	logStdout.Println("%s", lg)

	f, err := os.OpenFile("/var/log/artifakt/goapp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logStdout.Println("%s", err)
	}

	logFile := log.New(f, "goapp.log", log.LstdFlags)
	logFile.Println("%s", lg)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthCheckHandler)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/", homeHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
