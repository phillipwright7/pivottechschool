package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const endpoint string = "http://localhost:8080/products"

type product struct {
	Name        string `json:"name"`
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

func allProducts() {
	res, err := http.Get(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var products []product
	if err := json.Unmarshal(bs, &products); err != nil {
		log.Fatal(err)
	}
	for _, p := range products {
		fmt.Printf("%d: %s: $%d\n", p.ID, p.Name, p.Price)
	}
}

func oneProduct() {
	var id string
	var product []product

	fmt.Println("Which product ID would you like to return?")
	fmt.Scanln(&id)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("id", id)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &product)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req)
}

func main() {
	var input string
	fmt.Printf("Welcome to the product server. Please input the following numbers for whichever action you'd like to perform.\n1: View entire list of products.\n2: Return a specified product according to its ID.\n3: Add a new product to the product server.\n4: Update a product's name and price.\n5: Remove a product from the server.\n")
	fmt.Scanln(&input)
	switch input {
	case "1":
		allProducts()
	case "2":
		oneProduct()
	default:
		fmt.Println("Please input a valid number.")
	}
}
