package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		_, err := fmt.Fprintf(w, "ITS ALIVEEE LOOOLL!")
		if err != nil {
			fmt.Println("I dont even know this can error")
		}
	})

	server := &http.Server{
		Addr:         ":5500",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
