package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

var products []product

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatalf("error converting id to int: %v", err)
	}
	p := productFinder(id)
	if p == nil {
		log.Printf("product with id %d is not found", id)
	}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Printf("error encoding product: %v", id)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p product
	if err := json.NewEncoder(w).Encode(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	p.ID = len(products) + 1
	fmt.Println("Enter the name of the new product")
	fmt.Scanln(&p.Name)
	fmt.Println("Enter the price of the new product")
	fmt.Scanln(&p.Price)
	fmt.Println("Enter the description of the new product")
	fmt.Scanln(&p.Description)
	products = append(products, p)
	fmt.Println(products)
}

func productFinder(id int) *product {
	for _, p := range products {
		if p.ID == id {
			return &p
		}
	}
	return nil
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
	r.HandleFunc("/products/{id}", getProductHandler)
	r.HandleFunc("/products", addProductHandler).Methods(http.MethodPost)
	r.HandleFunc("/products", getProductsHandler)
	log.Println("Listening on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
