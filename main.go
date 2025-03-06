package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

// App holds application-wide dependencies, like the Cassandra session.
type App struct {
	DB            *gocql.Session
	RedisClient   *redis.Client
	KafkaProducer sarama.SyncProducer
}

func main() {

	// -------------------------------------------------------------------
	// 1. Cassandra Configuration
	// -------------------------------------------------------------------
	hosts := os.Getenv("CASSANDRA_HOSTS")
	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	user := strings.TrimSpace(os.Getenv("CASSANDRA_USERNAME"))
	pass := strings.TrimSpace(os.Getenv("CASSANDRA_PASSWORD"))

	cluster := gocql.NewCluster(strings.Split(hosts, ",")...)
	cluster.Keyspace = keyspace

	// If auth is required:
	if user != "" && pass != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{Username: user, Password: pass}
	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Error creating Cassandra session: %v", err)
	}
	defer session.Close()

	// -------------------------------------------------------------------
	// 2. Redis Configuration
	//    Reads REDIS_HOST, REDIS_PORT, REDIS_PASSWORD from environment
	// -------------------------------------------------------------------
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Build the address in "<host>:<port>" format
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// Create a Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	// Ping to test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Successfully connected to Redis")

	// -------------------------------------------------------------------
	// 3. Kafka Producer Configuration (using Shopify/sarama)
	// -------------------------------------------------------------------
	kafkaBrokers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS") // e.g. "localhost:9092"
	// Print out the values for debugging
	fmt.Printf("KAFKA_BOOTSTRAP_SERVERS: %s\n", kafkaBrokers)

	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}
	brokerList := strings.Split(kafkaBrokers, ",")

	// Setup Redis
	// Sarama config for a sync producer (simple example)
	saramaConfig := sarama.NewConfig()

	// Enable SASL (if that's how your Kafka is secured)
	saramaConfig.Net.SASL.Enable = true
	saramaConfig.Net.SASL.User = os.Getenv("KAFKA_SASL_USERNAME")
	saramaConfig.Net.SASL.Password = os.Getenv("KAFKA_SASL_PASSWORD")

	// Usually SASL mechanism is plain or scram-sha-256/512:
	saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, saramaConfig)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println("Failed to close Kafka producer:", err)
		}
	}()

	app := &App{DB: session}

	// Create a new router
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc("/accounts", app.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts", app.ReadAccount).Methods("GET")
	r.HandleFunc("/accounts/{accountId}", app.ReadAccounts).Methods("GET")
	r.HandleFunc("/accounts/{accountId}", app.UpdateAccount).Methods("PUT")
	r.HandleFunc("/accounts/{accountId}", app.DeleteAccount).Methods("DELETE")

	// Start the server on port 8080
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
