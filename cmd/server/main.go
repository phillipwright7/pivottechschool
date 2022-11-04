package main

import (
	"encoding/json"
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
		return
	}
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("error converting id to int: %v", err)
	}
	p := productFinder(id)
	if p == nil {
		log.Printf("product with id %d is not found", id)
		return
	}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Printf("error encoding product: %v", id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("error decoding product: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	p.ID = len(products) + 1
	products = append(products, p)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("error encoding product: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("error converting id to int: %v", err)
		return
	}
	p := productFinder(id)
	if p == nil {
		log.Printf("product with id %d is not found", id)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("error decoding product: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	products[id] = *p
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("error encoding product: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("error converting id to int: %v", err)
		return
	}
	p := productFinder(id)
	for i, prod := range products {
		if prod.ID == p.ID {
			products = append(products[:i], products[i+1:]...)
		}
	}
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("error encoding product: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	bs, err := os.ReadFile("cmd/server/products.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := json.Unmarshal(bs, &products); err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	initProducts()
	r := mux.NewRouter()
	r.HandleFunc("/products", addProductHandler).Methods(http.MethodPost)
	r.HandleFunc("/products/{id}", updateProductHandler).Methods(http.MethodPut)
	r.HandleFunc("/products/{id}", deleteProductHandler).Methods(http.MethodDelete)
	r.HandleFunc("/products/{id}", getProductHandler)
	r.HandleFunc("/products", getProductsHandler)
	log.Println("Listening on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
