package main

import "github.com/gorilla/mux"

func connectRoutes(r *mux.Router) {
	controller := newController(newRepository())
	r.HandleFunc("/", controller.hello).Methods("GET")
	r.HandleFunc("/transaction", controller.createTransaction).Methods("POST")
	r.HandleFunc("/transaction/{reference}", controller.fetchTransactionDetails).Methods("GET")
	r.HandleFunc("/account/{id}", controller.fetchUserAccountDetails).Methods("GET")
}
