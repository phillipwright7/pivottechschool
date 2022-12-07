package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

var database *sql.DB

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var payload []product
	params := r.URL.Query()["limit"]
	paramsConv, err := strconv.Atoi(params[0])
	if paramsConv < 0 {
		w.WriteHeader(400)
		return
	}
	if err != nil {
		w.WriteHeader(400)
		return
	}

	rows, err := database.Query("SELECT id, name, price FROM products ORDER BY id LIMIT ?", params[0])
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var price float64
		if err := rows.Scan(&id, &name, &price); err != nil {
			w.WriteHeader(500)
			return
		}
		p := product{
			ID:    id,
			Name:  name,
			Price: int(price),
		}
		payload = append(payload, p)
	}

	err = rows.Err()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	varConv, err := strconv.Atoi(vars["id"])
	if varConv < 0 {
		w.WriteHeader(400)
		return
	}
	if err != nil {
		w.WriteHeader(400)
		return
	}

	var payload product

	rows := database.QueryRow("SELECT id, name, price FROM products WHERE id = ?", vars["id"])
	var id int
	var name string
	var price float64
	if err := rows.Scan(&id, &name, &price); err != nil {
		w.WriteHeader(404)
		return
	}

	payload = product{
		ID:    id,
		Name:  name,
		Price: int(price),
	}

	err = rows.Err()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		w.WriteHeader(500)
		return
	}
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(400)
		return
	}

	p.Validate(w)

	tx, err := database.Begin()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	stmt, err := tx.Prepare("INSERT INTO products (name, price) values(?, ?)")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Name, p.Price)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	varConv, err := strconv.Atoi(vars["id"])
	if varConv < 0 {
		w.WriteHeader(400)
		return
	}
	if err != nil {
		w.WriteHeader(400)
		return
	}

	var p product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(500)
		return
	}
	p.Validate(w)

	rows := database.QueryRow("SELECT id, name, price FROM products WHERE id = ?", vars["id"])
	var id int
	var name string
	var price float64
	if err := rows.Scan(&id, &name, &price); err != nil {
		w.WriteHeader(404)
		return
	}

	tx, err := database.Begin()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	stmt, err := tx.Prepare("UPDATE products SET name = ?, price = ? WHERE id = ?")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(p.Name, p.Price, vars["id"]); err != nil {
		w.WriteHeader(500)
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(500)
		return
	}

}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	if _, err := strconv.Atoi(vars["id"]); err != nil {
		w.WriteHeader(400)
		return
	}

	rows := database.QueryRow("SELECT id, name, price FROM products WHERE id = ?", vars["id"])
	var id int
	var name string
	var price float64
	if err := rows.Scan(&id, &name, &price); err != nil {
		w.WriteHeader(404)
		return
	}

	tx, err := database.Begin()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	stmt, err := tx.Prepare("DELETE FROM products WHERE id = ?")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(vars["id"]); err != nil {
		w.WriteHeader(400)
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(500)
		return
	}
}

func initProducts() *sql.DB {
	var dbPtr string
	flag.StringVar(&dbPtr, "db", "products.db", "Flag to find products.db")
	flag.Parse()

	db, err := sql.Open("sqlite3", dbPtr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to products database successfully!")

	return db
}

func (p *product) Validate(w http.ResponseWriter) {
	if p.Name == "" || p.Price == 0 {
		w.WriteHeader(400)
		return
	}
}

func main() {
	database = initProducts()
	r := mux.NewRouter()
	r.HandleFunc("/products", addProductHandler).Methods(http.MethodPost)
	r.HandleFunc("/products/{id}", updateProductHandler).Methods(http.MethodPut)
	r.HandleFunc("/products/{id}", deleteProductHandler).Methods(http.MethodDelete)
	r.HandleFunc("/products/{id}", getProductHandler)
	r.HandleFunc("/products", getProductsHandler).Queries("limit", "{limit:.*}")
	log.Println("Listening on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
