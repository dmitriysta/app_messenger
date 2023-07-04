package messenger

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	os "os"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "qwerty123456"
	dbname   = "messenger"
)

func dbConnection() *sql.DB {
	os.Setenv("host", "localhost")
	os.Setenv("port", "5432")
	os.Setenv("user", "admin")
	os.Setenv("password", "qwerty123456")
	os.Setenv("dbname", "messenger")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return db
}

func createKafkaProducer() *kafka.Writer {
	config := kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "messages",
	}

	producer := kafka.NewWriter(config)

	return producer
}

func notifier() {
	config := kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "messages",
		GroupID:  "notifier-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	for {
		m, err := reader.ReadMessage()
		if err != nil {
			log.Printf("failed to read message: %v", err)
			continue
		}

		fmt.Printf("Received message: %s\n", m.Value)

	}

	handleNotificationMessage(m.Value)

}

func handleNotificationMessage(message []byte) {
	var notification GetMessageResponce
	err := json.Unmarshal(message, &notification)
	if err != nil {
		log.Printf("failed to unmarshal message: %v", err)
		return
	}

	log.Printf("Received notification: %+v", notification)
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "qwerty123456",
		DB:       1,
	})
	defer redisClient.Close()

	dbConnection()

	producer := createKafkaProducer()
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	go notifier()

	r := mux.NewRouter()

	messageController := NewMessageController(logger)

	r.HandleFunc("/messages", messageController.CreateMessage).Methods("POST")
	r.HandleFunc("/messages/{id}", messageController.ReceiveMessage).Methods("GET")
	r.HandleFunc("/messages/{id}", messageController.SendMessage).Methods("PUT")
	r.HandleFunc("/messages/{id}", messageController.DeleteMessage).Methods("DELETE")

	logger.Info("successfully connected to PostgreSQL database")

	port := ":8080"
	logger.Info("server listening on port", zap.String("port", port))

	err = http.ListenAndServe(port, r)
	if err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
