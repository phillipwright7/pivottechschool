# Capstone Project

## Summary

The purpose of this capstone project is to create a shopping cart that will retrieve and parse product JSON data from the site https://dummyjson.com/products, prompt the user to pick a category of product, and then the user will have the ability to pick a product from the selected category and add that product to a cart which is stored in memory.

I will demonstrate the following:

1. Utilizing HTTP GET methods in order to access remotely stored data
2. Unmarshalling JSON into data structures native to Golang
3. Prompting the user via the CLI and allowing them to make inputs
4. Writing a program in Golang that will extract the user's inputs and run conditional actions based on what the user wants
5. Planning for and handling errors in Golang
6. Presenting the product data in a readable format
7. Appending selected items to an in-memory shopping cart that the user can view

## User Stories

### As a user, I would to view every product category, and then select a product to view a list of relevant products to purchase.

**Acceptance Criteria**

When the program is run, a list of product categories will be shown, each numbered and on a new line. At the bottom, the user will see a prompt to enter a number or "q" to quit the program, then press enter to finish.

Example:
```
$ go run main.go
```

```
Select a product category:
1. ...
2. ...
3. ...
4. ...

Input a number or "q" to quit, then press enter.
```

### As a user, I would like to see a list of items and be able to easily select any number of them to add to my cart

**Acceptance Criteria**

When the user is prompted with a list of products under a certain cateogry, each will be numbered and on a new line. At the bottom, the user will see a prompt to input any amount of numbers, each seperated by a space, and to press enter to finish. Or, they can simple enter "q" and press enter to quit.

Example:
```
Select a product:
1. ...
2. ...
3. ...
4. ...

Input any numbers you wish to add to your cart, each seperated by a space, then press enter.
(Or if you would like to quit, press "q" and then press enter.)
```

### As a user, I would like to view my cart and be able to delete or edit the quantity of each item.

**Acceptance Criteria**

At any point after the program has started, the user will be able to enter "c" and then press enter to view their cart. The program will list each item in the order it was added, and have the quantity listed beside it. The user will be able to select a number of the product they want to change the quantity for, or they can enter "r" to return to the categories view.

Example:

```
$ go run main.go
```

```
$ c
```

```
Your cart:
1. ... | price: $... | quantity: #
2. ... | price: $... | quantity: #
3. ... | price: $... | quantity: #
4. ... | price: $... | quantity: #

Input a product number for the quantity you'd like to edit.
(Or if you'd like to return to the category view, press "r" and then press enter.)
```

```
Edit quantity:
3. ... (qty. # --> ?)

Input the quanity number you'd like to change for ...
(A value of 0 will delete the item from your cart.)
```