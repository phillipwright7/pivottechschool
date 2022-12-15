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

var cartTitles = []string{}

var cartCount, cartPrice, cartQuantity = map[string]int{}, map[string]int{}, map[string]int{}

var countPad, titlePad, brandPad, pricePad, quantityPad int

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
		return input, nil
	}

	if input == "c" {
		return input, nil
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

func productList(c string) ([]product, error) {
	count := 1
	var products []product

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
	fmt.Println(strings.Repeat("-", len(fmt.Sprint(count))+titlePad+brandPad+pricePad+9))

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

func shoppingCart() (string, int, error) {
	var total float64

	for _, p := range shoppingList {
		if _, ok := cartCount[p.Title]; !ok {
			cartCount[p.Title] = len(cartCount) + 1
			cartTitles = append(cartTitles, p.Title)
		}
		if _, ok := cartPrice[p.Title]; !ok {
			cartPrice[p.Title] = p.Price
		}
		if _, ok := cartQuantity[p.Title]; !ok {
			cartQuantity[p.Title] = 0
		}
		if _, ok := cartQuantity[p.Title]; ok {
			cartQuantity[p.Title]++
		}
	}

	for i := 0; i < len(cartCount); i++ {
		cartCount[cartTitles[i]] = i + 1
	}

	fmt.Printf("\nYour shopping cart:\n")

	for _, p := range shoppingList {
		if strings.Count(strconv.Itoa(cartCount[p.Title]), "") > countPad {
			countPad = strings.Count(strconv.Itoa(cartCount[p.Title]), "")
		}
		if strings.Count(p.Title, "") > titlePad {
			titlePad = strings.Count(p.Title, "")
		}
		if strings.Count(strconv.Itoa(p.Price), "") > pricePad {
			pricePad = strings.Count(strconv.Itoa(p.Price), "")
		}
		if strings.Count(fmt.Sprint(cartQuantity[p.Title]), "") > quantityPad {
			quantityPad = strings.Count(fmt.Sprint(cartQuantity[p.Title]), "")
		}
	}

	fmt.Printf("%*s %-*s   %-*s   %s\n", countPad, " ", titlePad, "Name", pricePad, "Price", "Quantity")
	fmt.Println(strings.Repeat("-", countPad+titlePad+pricePad+quantityPad+9))

	for _, t := range cartTitles {
		fmt.Printf("%0*d. %-*s | $%-*d | %-*d\n", countPad, cartCount[t], titlePad, t, pricePad, cartPrice[t], quantityPad, cartQuantity[t])
		total += float64(cartPrice[t] * cartQuantity[t])
	}

	fmt.Println(strings.Repeat("-", countPad+titlePad+pricePad+quantityPad+9))
	fmt.Printf("Total: $%.2f\n", total)
	fmt.Printf("\nInput a product number for the quantity you'd like to edit.\n(Or if you'd like to return to the category view, press 'r' and then press enter.\n")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", 0, err
	}
	input = strings.Replace(input, "\n", "", -1)

	if input == "r" {
		return input, 0, nil
	}

	inputNum, err := strconv.Atoi(input)
	if err != nil {
		return "", 0, errors.New("error: not a number")
	}
	if inputNum > len(cartQuantity) || inputNum < 1 {
		return "", 0, errors.New("error: not a valid number")
	}

	return input, inputNum, nil
}

func editCart(i int) error {
	index := i - 1
	fmt.Printf("\nEdit quantity:\n")
	fmt.Printf("%s (qty. %d --> ?)\n", cartTitles[index], cartQuantity[cartTitles[index]])
	fmt.Printf("Input the quanity number you'd like to change for %s\n(A value of 0 will delete the item from your cart.)\n", cartTitles[index])

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.Replace(input, "\n", "", -1)

	inputNum, err := strconv.Atoi(input)
	if err != nil {
		return errors.New("error: not a number")
	}
	if inputNum == 0 {
		delete(cartCount, cartTitles[index])
		delete(cartPrice, cartTitles[index])
		delete(cartQuantity, cartTitles[index])
		cartTitles = append(cartTitles[:index], cartTitles[index+1:]...)
		return nil
	}

	unitPrice := cartPrice[cartTitles[index]] / cartQuantity[cartTitles[index]]
	cartQuantity[cartTitles[index]] = inputNum
	cartPrice[cartTitles[index]] = unitPrice * inputNum

	return nil
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
	if category == "q" {
		log.Println("program has quit")
		return
	}
CART:
	if category == "c" && err == nil {
		input, editNum, err := shoppingCart()
		if err != nil {
			log.Println(err)
			goto CART
		}
		if input == "r" {
			shoppingList = []product{}
			goto CATEGORIES
		}
		if err = editCart(editNum); err != nil {
			log.Println(err)
			goto CART
		}
		if err == nil {
			shoppingList = []product{}
			goto CART
		}
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
