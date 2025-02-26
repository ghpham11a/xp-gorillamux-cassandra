package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

func main() {
    // Create a new router
    r := mux.NewRouter()

    // Register a simple handler for the root path
    r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        fmt.Fprintln(w, "Hello, World!")
    })

    // Start the server on port 8080
    log.Println("Starting server on :8080...")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}
