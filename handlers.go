package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func (app *App) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse JSON into a struct
	var account Account
	if err := json.Unmarshal(reqBody, &account); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Respond with the parsed struct (for example purposes)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		log.Println("Error writing response:", err)
	}
}

func (app *App) ReadAccounts(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: Readaccounts")

	// Create a list of 2 accounts
	mockAccounts := []Account{
		{
			ID:            gocql.TimeUUID(),
			Email:         "john.doe@example.com",
			DateOfBirth:   time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			AccountNumber: "ACC-12345",
			Balance:       100.50,
			CreatedAt:     time.Now(),
		},
		{
			ID:            gocql.TimeUUID(),
			Email:         "jane.smith@example.com",
			DateOfBirth:   time.Date(1985, 12, 5, 0, 0, 0, 0, time.UTC),
			AccountNumber: "ACC-67890",
			Balance:       250.00,
			CreatedAt:     time.Now(),
		},
	}

	// Respond with the list of accounts
	if err := json.NewEncoder(w).Encode(mockAccounts); err != nil {
		http.Error(w, "Failed to encode accounts: "+err.Error(), http.StatusInternalServerError)
	}

	// fmt.Println("Endpoint hit: returnAllGroceries")
	// json.NewEncoder(w).Encode(groceries)
}

// type Account struct {
// 	ID            gocql.UUID `json:"id"` // PRIMARY KEY in Cassandra
// 	Email         string     `json:"email"`
// 	DateOfBirth   time.Time  `json:"date_of_birth"`
// 	AccountNumber string     `json:"account_number"`
// 	Balance       float64    `json:"balance"`
// 	CreatedAt     time.Time  `json:"created_at"`
// }

func (app *App) ReadAccount(w http.ResponseWriter, r *http.Request) {
	// 1. Get the path parameter (e.g. /accounts/{accountId})
	vars := mux.Vars(r)
	accountId := vars["accountId"]
	if accountId == "" {
		http.Error(w, "accountId path parameter is required", http.StatusBadRequest)
		return
	}

	// --------------------------------------
	// Try Redis first
	// --------------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	redisKey := fmt.Sprintf("account:%d", accountId)
	cachedVal, cacheErr := app.RedisClient.Get(ctx, redisKey).Result()
	if cacheErr == nil {
		// Cache hit
		var cachedAccount Account
		if err := json.Unmarshal([]byte(cachedVal), &cachedAccount); err == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cachedAccount)
			return
		}
		// If we fail to unmarshal, proceed to DB fetch
	} else if cacheErr != redis.Nil {
		// Some other Redis error (network, etc.)
		log.Println("Warning: error reading from Redis:", cacheErr)
	}

	// 2. Define a struct to hold the result
	var account Account
	// 3. Build and execute a Cassandra query
	//    Adjust the SELECT fields/WHERE clause to match your Cassandra schema.
	queryErr := app.DB.Query(
		"SELECT ID, Email, DateOfBirth, AccountNumber, Balance, CreatedAt FROM accounts WHERE ID = ? LIMIT 1",
		accountId,
	).Consistency(gocql.One).Scan(&account.ID, &account.Email, &account.DateOfBirth, &account.AccountNumber, &account.Balance, &account.CreatedAt)

	if queryErr != nil {
		if queryErr == gocql.ErrNotFound {
			// 4. No record found for that ID
			http.NotFound(w, r)
			return
		}
		// 5. Some other error occurred
		http.Error(w, "Error querying Cassandra: "+queryErr.Error(), http.StatusInternalServerError)
		return
	}

	// --------------------------------------
	// Store in Redis for next time
	// --------------------------------------
	empBytes, _ := json.Marshal(account)
	if err := app.RedisClient.Set(ctx, redisKey, empBytes, 5*time.Minute).Err(); err != nil {
		log.Println("Warning: failed to cache employee in Redis:", err)
	}

	// 6. On success, encode the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (app *App) UpdateAccount(w http.ResponseWriter, r *http.Request) {

}

func (app *App) DeleteAccount(w http.ResponseWriter, r *http.Request) {

}
