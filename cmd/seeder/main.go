package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type products []product

func main() {
	var dbPtr string
	var jsonPtr string
	flag.StringVar(&dbPtr, "d", "products.db", "Flag to find products.db")
	flag.StringVar(&jsonPtr, "j", "products.json", "Flag to read products.json")
	flag.Parse()

	if err := os.RemoveAll(dbPtr); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbPtr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := "CREATE TABLE products (id INTEGER NOT NULL PRIMARY KEY, name TEXT, price REAL)"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	prod, err := os.ReadFile(jsonPtr)
	if err != nil {
		log.Fatal(err)
	}
	var payload products
	if err := json.Unmarshal(prod, &payload); err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO products (id, name, price) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, p := range payload {
		_, err = stmt.Exec(p.ID, p.Name, p.Price)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, name, price FROM products WHERE id <= 5")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var price float64
		if err := rows.Scan(&id, &name, &price); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, price)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
