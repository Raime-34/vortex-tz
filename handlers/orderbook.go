package handlers

import (
	"encoding/json"
	"golang.org/x/text/encoding/charmap"
	"net/http"
	"vortex-tz/models"
	"vortex-tz/queries"
)

func SaveOrderHandle(w http.ResponseWriter, r *http.Request) {
	var orderBook models.OrderBook
	err := json.NewDecoder(r.Body).Decode(&orderBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = queries.SaveOrderBook(orderBook.Exchange, orderBook.Pair, orderBook.Asks, orderBook.Bids)
	if err != nil {
		encoder := charmap.Windows1251.NewDecoder()
		s, e := encoder.String(err.Error())
		if e != nil {
			println(e.Error())
		}
		http.Error(w, s, http.StatusBadRequest)
	}
}

func GetOrderHandle(w http.ResponseWriter, r *http.Request) {
	exchangeName := r.URL.Query().Get("exchange")
	pair := r.URL.Query().Get("pair")

	orderBook, err := queries.GetOrderBook(exchangeName, pair)
	if err != nil {
		encoder := charmap.Windows1251.NewDecoder()
		s, e := encoder.String(err.Error())
		if e != nil {
			println(e.Error())
		}
		http.Error(w, s, http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(orderBook)
	if err != nil {
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}
	w.Write(bytes)
}
