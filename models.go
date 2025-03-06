package main

import (
	"time"

	"github.com/gocql/gocql"
)

// Account represents the "accounts" table structure in Cassandra.
type Account struct {
	ID            gocql.UUID `json:"id"` // PRIMARY KEY in Cassandra
	Email         string     `json:"email"`
	DateOfBirth   time.Time  `json:"date_of_birth"`
	AccountNumber string     `json:"account_number"`
	Balance       float64    `json:"balance"`
	CreatedAt     time.Time  `json:"created_at"`
}
