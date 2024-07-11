package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/midedickson/simple-banking-app/config"
	"github.com/midedickson/simple-banking-app/routes"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	config.ConnectToDB()
	config.AutoMigrate()
	r := mux.NewRouter()
	routes.ConnectRoutes(r)
	log.Println("Starting Simple Banking Server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
