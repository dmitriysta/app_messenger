package messenger

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "qwerty123456"
	dbname   = "messenger"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatal("Failed to open database connection: %v", zap.Error(err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	r := mux.NewRouter()
	userController := NewUserController(logger)
	r.HandleFunc("/users", userController.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", userController.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id}", userController.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")

	logger.Info("successfully connected to PostgreSQL database")

	port := ":80"
	logger.Info("server listening on port", zap.String("port", port))

	err = http.ListenAndServe(port, r)
	if err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}

}
