package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"vortex-tz/handlers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/orderbook/", handlers.SaveOrderHandle).Methods("POST")
	router.HandleFunc("/orderbook", handlers.GetOrderHandle).Methods("GET")

	http.ListenAndServe(":8080", router)
}

//func GetOrderBook(exchangeName string, pair string) ([]*DepthOrder, error) {
//
//}
