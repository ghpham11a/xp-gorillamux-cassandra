package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

// App holds application-wide dependencies, like the Cassandra session.
type App struct {
	DB *gocql.Session
}

func main() {

	hosts := os.Getenv("CASSANDRA_HOSTS")
	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	user := strings.TrimSpace(os.Getenv("CASSANDRA_USERNAME"))
	pass := strings.TrimSpace(os.Getenv("CASSANDRA_PASSWORD"))

	// Print out the values for debugging
	fmt.Printf("CASSANDRA_HOSTS: %s\n", hosts)
	fmt.Printf("CASSANDRA_KEYSPACE: %s\n", keyspace)
	fmt.Printf("CASSANDRA_USERNAME: %s\n", user)
	fmt.Printf("CASSANDRA_PASSWORD: %s\n", pass)

	cluster := gocql.NewCluster(strings.Split(hosts, ",")...)
	cluster.Keyspace = keyspace

	// If auth is required:
	if user != "" && pass != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{Username: user, Password: pass}
	}

	// cluster.Authenticator = gocql.PasswordAuthenticator{Username: "cassandra", Password: "oRRn8ZbkB0"}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Error creating Cassandra session: %v", err)
	}
	defer session.Close()

	app := &App{DB: session}

	// Create a new router
	r := mux.NewRouter()

	r.HandleFunc("/employees", app.CreateEmployee).Methods("POST")
	r.HandleFunc("/employees", app.ReadEmployees).Methods("GET")
	r.HandleFunc("/employees/{empId}", app.ReadEmployee).Methods("GET")
	r.HandleFunc("/employees/{empId}", app.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employees/{empId}", app.DeleteEmployee).Methods("DELETE")

	// Start the server on port 8080
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
