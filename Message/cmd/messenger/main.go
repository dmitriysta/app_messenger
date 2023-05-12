package messenger

import (
	"database/sql"
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
