package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("connection received")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	bind := os.Getenv("BIND_ADDRESS")

	if bind == "" {
		bind = "127.0.0.1"
	}

	server := &http.Server{
		Addr: fmt.Sprintf("%s:8000", bind),
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	
	log.Fatal(server.ListenAndServe())
}
