package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type response struct {
	Products []product `json:"products"`
	Total    int       `json:"total"`
	Skip     int       `json:"skip"`
	Limit    int       `json:"limit"`
}

type product struct {
	ID                 int      `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Price              int      `json:"price"`
	DiscountPercentage float64  `json:"discountPercentage"`
	Rating             float64  `json:"rating"`
	Stock              int      `json:"stock"`
	Brand              string   `json:"brand"`
	Category           string   `json:"category"`
	Thumbnail          string   `json:"thumbnail"`
	Images             []string `json:"images"`
}

var prodResp response

var shoppingList []product

func categoryList() (string, error) {
	cat := map[int]string{}

	for _, c := range prodResp.Products {
		if len(cat) == 0 {
			cat[1] = c.Category
		}
		if cat[len(cat)] != c.Category {
			cat[len(cat)+1] = c.Category
		}
	}
	fmt.Println("Select a product category:")
	for i := 1; i <= len(cat); i++ {
		fmt.Println(fmt.Sprint(i) + ": " + cat[i])
	}
	fmt.Printf("\nInput a number or 'q' to quit, then press enter.\n(Or if you would like to view your cart, press 'c' and then press enter.\n")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.Replace(input, "\n", "", -1)

	if input == "q" {
		log.Fatalln("program has quit")
	}

	if input == "c" {
		return "", nil
	}

	resp, err := strconv.Atoi(input)
	if err != nil {
		return "", errors.New("error: invalid number")
	}
	if _, ok := cat[resp]; !ok {
		return "", errors.New("error: category doesn't exist")
	}

	fmt.Printf("\n")
	return cat[resp], nil
}

func shoppingCart() {
	count := 1
	var total float64
	var titlePad, brandPad, pricePad int

	fmt.Printf("\nYour shopping cart:\n")

	for _, p := range shoppingList {
		if len(p.Title) > titlePad {
			titlePad = len(p.Title)
		}
		if len(p.Brand) > brandPad {
			brandPad = len(p.Brand)
		}
		if len(fmt.Sprint(p.Price)) > pricePad {
			pricePad = len(fmt.Sprint(p.Price))
		}
	}

	fmt.Printf("%*s %-*s   %-*s   %s\n", len(fmt.Sprint(count))+1, " ", titlePad, "Name", brandPad, "Brand", "Price")
	fmt.Println(strings.Repeat("-", len(fmt.Sprint(count))+titlePad+brandPad+pricePad+9))
	for _, p := range shoppingList {
		fmt.Printf("%s: %-*s | ", fmt.Sprint(count), titlePad, p.Title)
		fmt.Printf("%-*s | ", brandPad, p.Brand)
		fmt.Printf("$%s\n", fmt.Sprint(p.Price))
		count++
		total += float64(p.Price)
	}

	fmt.Println(strings.Repeat("-", len(fmt.Sprint(count))+titlePad+brandPad+pricePad+9))
	fmt.Printf("Total: %.2f\n", total)
	fmt.Printf("\nInput a product number for the quantity you'd like to edit.\n(Or if you'd like to return to the category view, press 'r' and then press enter.\n")
}

func productList(c string) ([]product, error) {
	count := 1
	var products []product
	var titlePad, brandPad, pricePad int

	fmt.Println("Select a product:")
	for _, p := range prodResp.Products {
		if p.Category == c {
			if len(p.Title) > titlePad {
				titlePad = len(p.Title)
			}
			if len(p.Brand) > brandPad {
				brandPad = len(p.Brand)
			}
			if len(fmt.Sprint(p.Price)) > pricePad {
				pricePad = len(fmt.Sprint(p.Price))
			}
		}
	}

	fmt.Printf("%*s %-*s   %-*s   %s\n", len(fmt.Sprint(count))+1, " ", titlePad, "Name", brandPad, "Brand", "Price")
	fmt.Println(strings.Repeat("-", len(fmt.Sprint(count))+titlePad+brandPad+pricePad+9))
	for _, p := range prodResp.Products {
		if p.Category == c {
			products = append(products, p)
			fmt.Printf("%s: %-*s | ", fmt.Sprint(count), titlePad, p.Title)
			fmt.Printf("%-*s | ", brandPad, p.Brand)
			fmt.Printf("$%s\n", fmt.Sprint(p.Price))
			count++
		}
	}

	fmt.Printf("\nInput any numbers you wish to add to your cart, each seperated by a space, then press enter.\n(Or if you would like to return to category view, press 'r' and then press enter.)\n")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.Replace(input, "\n", "", -1)
	inputSlice := strings.Split(input, " ")

	if input == "r" {
		return shoppingList, nil
	}

	for _, s := range inputSlice {
		i, err := strconv.Atoi(s)
		if err != nil {
			return []product{}, errors.New("error: not a number")
		}
		if i > len(products) || i < 1 {
			return []product{}, errors.New("error: invalid number")
		}
	}

	for _, n := range inputSlice {
		for i, p := range products {
			if n == fmt.Sprint(i+1) {
				shoppingList = append(shoppingList, p)
				fmt.Println(p.Title, "added to your cart!")
			}
		}
	}
	return shoppingList, nil
}

func main() {
	resp, err := http.Get("https://dummyjson.com/products")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &prodResp); err != nil {
		log.Fatal(err)
	}

CATEGORIES:
	category, err := categoryList()
	if err != nil {
		log.Println(err)
		fmt.Printf("\n")
		goto CATEGORIES
	}
	if category == "" && err == nil {
		shoppingCart()
		return
	}

PRODUCTS:
	shoppingList, err := productList(category)
	if err != nil {
		log.Println(err)
		fmt.Printf("\n")
		goto PRODUCTS
	}
	if shoppingList == nil {
		fmt.Printf("\n")
		goto CATEGORIES
	}
	if shoppingList != nil && err == nil {
		fmt.Printf("\n")
		goto CATEGORIES
	}
}
