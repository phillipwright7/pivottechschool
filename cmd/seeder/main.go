package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var prodPth string
	flag.StringVar(&prodPth, "p", "products.db", "Flag to find products.db")
	flag.Parse()

	if err := os.RemoveAll(prodPth); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", prodPth)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := "CREATE TABLE products (id INTEGER NOT NULL PRIMARY KEY, name TEXT, price REAL)"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	prod, err := os.ReadFile("../server/products.json")
	if err != nil {
		log.Fatal(err)
	}
	var payload []map[string]interface{}
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

	for i := range payload {
		_, err = stmt.Exec(payload[i]["id"], payload[i]["name"], payload[i]["price"])
		if err != nil {
			log.Fatal(err)
		}
	}

	rows, err := db.Query("SELECT id, name, price FROM products WHERE id = 1")
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
	spew.Dump(rows)
}
