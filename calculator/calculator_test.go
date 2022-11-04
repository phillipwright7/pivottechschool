package calculator_test

import (
	"testing"

	"github.com/phillipwright7/pivottechschool/calculator"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name                string
		input, input2, want int
	}{
		{name: "first", input: 2, input2: 2, want: 4},
		{name: "second", input: 2, input2: 3, want: 5},
		{name: "third", input: 0, input2: 2, want: 2},
		{name: "fourth", input: -1, input2: 2, want: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := calculator.Add(test.input, test.input2)
			if got != test.want {
				t.Errorf("got: %q, want %q", got, test.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name                string
		input, input2, want int
	}{
		{name: "first", input: 2, input2: 2, want: 0},
		{name: "second", input: 2, input2: 3, want: -1},
		{name: "third", input: 0, input2: 2, want: -2},
		{name: "fourth", input: -1, input2: 2, want: -3},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := calculator.Subtract(test.input, test.input2)
			if got != test.want {
				t.Errorf("got: %q, want %q", got, test.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name                string
		input, input2, want int
	}{
		{name: "first", input: 2, input2: 2, want: 4},
		{name: "second", input: 2, input2: 3, want: 6},
		{name: "third", input: 0, input2: 2, want: 0},
		{name: "fourth", input: -1, input2: 2, want: -2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := calculator.Multiply(test.input, test.input2)
			if got != test.want {
				t.Errorf("got: %q, want %q", got, test.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name                string
		input, input2, want int
	}{
		{name: "first", input: 4, input2: 2, want: 2},
		{name: "second", input: 8, input2: 2, want: 4},
		{name: "third", input: 16, input2: 4, want: 4},
		{name: "fourth", input: 33, input2: 11, want: 3},
		{name: "divideByZero", input: 2, input2: 0, want: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := calculator.Divide(test.input, test.input2)
			if got != test.want && err != nil {
				t.Errorf("got: %q, want %q", got, test.want)
			}
		})
	}
}
