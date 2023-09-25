package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
	log.Printf("Received request for /hello from %s\n", req.RemoteAddr)
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	log.Printf("Received request for /headers from %s\n", req.RemoteAddr)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090" // Default port if PORT is not set
	}

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	log.Printf("Server is running on port %s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
