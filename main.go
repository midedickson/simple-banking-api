package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/midedickson/simple-banking-app/config"
	"github.com/midedickson/simple-banking-app/controllers"
	"github.com/midedickson/simple-banking-app/external"
	mock_client "github.com/midedickson/simple-banking-app/mock"
	"github.com/midedickson/simple-banking-app/repository"
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
	storageRepository := repository.NewStorageRepository(config.DB)
	mockClient := mock_client.CreateNewPOSTMockClient()
	external := external.NewTransactionExternal(mockClient)
	controller := controllers.NewController(storageRepository, external)
	routes.ConnectRoutes(r, controller)
	log.Println("Starting Simple Banking Server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
