package main

import (
	"fmt"

	"github.com/phillipwright7/pivottechschool/cmd/calculator"
)

func main() {
	a := 2
	b := 2
	sum := calculator.Add(a, b)
	fmt.Println(sum)
	dif := calculator.Subtract(a, b)
	fmt.Println(dif)
	pro := calculator.Multiply(a, b)
	fmt.Println(pro)
	quo, err := calculator.Divide(a, b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(quo)
	}
}
