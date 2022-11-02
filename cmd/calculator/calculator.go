package calculator

import "errors"

func Add(a, b int) int {
	sum := a + b
	return sum
}

func Subtract(a, b int) int {
	dif := a - b
	return dif
}

func Multiply(a, b int) int {
	pro := a * b
	return pro
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("divide by zero error")
	} else {
		quo := a / b
		return quo, nil
	}
}
