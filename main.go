package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())

	s := r.PathPrefix("/v0").Subrouter().StrictSlash(true)
	s.NotFoundHandler = http.HandlerFunc(notFound)
	s.HandleFunc("/new-sheet", streamHandler(prepareSheet)).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", logHandler(os.Stdout, r)))
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	log.Printf("Error status code: 404 when serving path: %s",
		r.RequestURI)
}
