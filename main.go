package main

import (
	"log"
	"net/http"
)

func launchWebServer(mux *http.ServeMux) {
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func main() {
	// Register the two new handler functions and corresponding URL patterns with
	// the servemux, in exactly the same way that we did before.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	go launchWebServer(mux)
	select {}
}
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from home page\n"))
}
