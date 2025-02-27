package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func (app *App) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse JSON into a struct
	var employee Employee
	if err := json.Unmarshal(reqBody, &employee); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Respond with the parsed struct (for example purposes)
	if err := json.NewEncoder(w).Encode(employee); err != nil {
		log.Println("Error writing response:", err)
	}
}

func (app *App) ReadEmployees(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: ReadEmployees")

	// Create a list of 2 employees
	mockEmployees := []Employee{
		{EmpID: 1, DeptID: 10, FirstName: "John", LastName: "Smith"},
		{EmpID: 2, DeptID: 20, FirstName: "Jane", LastName: "Doe"},
	}

	// Respond with the list of employees
	if err := json.NewEncoder(w).Encode(mockEmployees); err != nil {
		http.Error(w, "Failed to encode employees: "+err.Error(), http.StatusInternalServerError)
	}

	// fmt.Println("Endpoint hit: returnAllGroceries")
	// json.NewEncoder(w).Encode(groceries)
}

func (app *App) ReadEmployee(w http.ResponseWriter, r *http.Request) {
	// 1. Get the path parameter (e.g. /employees/{empId})
	vars := mux.Vars(r)
	empID := vars["empId"]
	if empID == "" {
		http.Error(w, "empId path parameter is required", http.StatusBadRequest)
		return
	}

	// 2. Define a struct to hold the result
	var employee Employee

	// 3. Build and execute a Cassandra query
	//    Adjust the SELECT fields/WHERE clause to match your Cassandra schema.
	err := app.DB.Query(
		"SELECT empID, deptID, first_name, last_name FROM employees WHERE empID = ? LIMIT 1",
		empID,
	).Consistency(gocql.One).Scan(&employee.EmpID, &employee.DeptID, &employee.FirstName, &employee.LastName)

	if err != nil {
		if err == gocql.ErrNotFound {
			// 4. No record found for that ID
			http.NotFound(w, r)
			return
		}
		// 5. Some other error occurred
		http.Error(w, "Error querying Cassandra: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. On success, encode the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

func (app *App) UpdateEmployee(w http.ResponseWriter, r *http.Request) {

}

func (app *App) DeleteEmployee(w http.ResponseWriter, r *http.Request) {

}
