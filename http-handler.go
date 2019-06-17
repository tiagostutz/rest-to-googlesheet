package main

import (
	"log"
	"net/http"
	"time"
)

func streamHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("Logged connection from: ", r.RemoteAddr)
		log.Println("starting at: ", time.Now())

		next.ServeHTTP(w, r)

		log.Println("finished at: ", time.Now())
	}
}
