package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	connectToDB()
	autoMigrate()
	r := mux.NewRouter()
	connectRoutes(r)
	log.Println("Starting Simple Banking Server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
