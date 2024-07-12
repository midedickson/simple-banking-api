package routes

import (
	"github.com/gorilla/mux"
	"github.com/midedickson/simple-banking-app/controllers"
)

func ConnectRoutes(r *mux.Router, controller *controllers.Controller) {
	r.HandleFunc("/", controller.Hello).Methods("GET")
	r.HandleFunc("/transaction", controller.CreateTransaction).Methods("POST")
	r.HandleFunc("/transaction/{reference}", controller.FetchTransactionDetails).Methods("GET")
	r.HandleFunc("/account/{id}", controller.FetchUserAccountDetails).Methods("GET")
}
