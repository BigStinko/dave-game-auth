package main

import (
	"log"
	"net/http"
)

const (
	LOCALADDR =   "127.0.0.1"
	DEFAULTPORT = "54321"
)

func main() {
	address, port := getAddressAndPort()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", getReadinessHandler)
	mux.HandleFunc("GET /err", getErrHandler)

	//tlsConfig, err := getTLSConfig()
	//if err != nil { panic(err) }

	server := &http.Server{
		Addr: address + ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

// For now just gives localhost
func getAddressAndPort() (address, port string) {
	return LOCALADDR, DEFAULTPORT
}

