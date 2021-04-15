package main

import (
	"database/sql"
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
	generateLog(r.Host)

	mysql, m := checkConnectionDB()
	msg = msg + "Check MYSQL Connection => " + mysql + "\n"
	msg += msg + m
	io.WriteString(w, msg)

}

func checkConnectionDB() (string, string) {
	con := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("ARTIFAKT_MYSQL_USER"), os.Getenv("ARTIFAKT_MYSQL_PASSWORD"), os.Getenv("ARTIFAKT_MYSQL_HOST"), os.Getenv("ARTIFAKT_MYSQL_PORT"), os.Getenv("ARTIFAKT_MYSQL_DATABASE_NAME"))
	db, err := sql.Open("mysql", con)
	defer db.Close()
	if err != nil {
		return "ko", con
	}
	return "ok", con
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
