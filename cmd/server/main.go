package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type product struct {
	Name        string `json:"name"`
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

const port string = ":8080"

var products []product

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func initProducts() {
	bs, err := os.ReadFile("products.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(bs, &products); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initProducts()

	r := mux.NewRouter()
	r.HandleFunc("/products", getProductsHandler)

	log.Printf("Listening on port %s...", port)
	log.Fatal(http.ListenAndServe(port, r))
}
