package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type product struct {
	Name        string `json:"name"`
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

func productsGet() []product {
	res, err := http.Get("http://localhost:8080/products")
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
	return products
}

func allProducts() {
	products := productsGet()
	for _, p := range products {
		fmt.Printf("%d: %s: $%d\n", p.ID, p.Name, p.Price)
	}
}

func oneProduct() {
	var id int64
	fmt.Println("Which product ID would you like to return?")
	fmt.Scanln(&id)
	products := productsGet()
	for i, p := range products {
		if int64(i) == id {
			fmt.Printf("%d: %s: $%d\n", p.ID, p.Name, p.Price)
		}
	}
}

func addProduct() {
	var name, description string
	var price int64

	fmt.Println("Please enter the product name:")
	fmt.Scanln(&name)
	fmt.Println("Please enter the product description:")
	fmt.Scanln(&description)
	fmt.Println("Please enter the product price:")
	fmt.Scanln(&price)

	newProd := product{
		Name:        name,
		ID:          101,
		Description: description,
		Price:       price,
	}

	reqBody, err := json.Marshal(newProd)
	res, err := http.Post("http://localhost:8080/products", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	fmt.Printf("%s added to product server.", name)
}

func main() {
	var input string
	fmt.Println("Welcome to the product server. Please input the following numbers for whichever action you'd like to perform.")
	fmt.Println("1: View entire list of products.")
	fmt.Println("2: Return a specified product according to its ID.")
	fmt.Println("3: Add a new product to the product server.")
	fmt.Println("4: Update a product's name and price.")
	fmt.Println("5: Remove a product from the server")
	fmt.Scanln(&input)
	switch input {
	case "1":
		allProducts()
	case "2":
		oneProduct()
	case "3":
		addProduct()
	default:
		fmt.Println("Please input a valid number.")
	}
}
