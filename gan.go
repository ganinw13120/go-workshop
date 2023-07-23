package main

import (
	"errors"
	"fmt"
	"sync"
)

type IPromotion interface {
	CalculatePrice(int) int
}

type Promotion struct {
	discount int
}

func (p Promotion) CalculatePrice(price int) int {
	return price - p.discount
}

type SpecialPromotion struct {
	discount int
}

func (p SpecialPromotion) CalculatePrice(price int) int {
	return price / p.discount
}

func generatePromotion() IPromotion {
	return SpecialPromotion{10}
}

type customError struct {
	message string
}

func (e customError) Error() string {
	return fmt.Sprintf("Error is %s", e.message)
}

func newError(message string) error {
	return customError{
		message: message,
	}
}

func main() {
	var err error
	err = customError{"hello"}
	err = nil
	err = newError("Not found!")
	err = errors.New("Hello")
	err.Error()

	a := 0

	var t sync.Mutex // Cannot be safely copy
	go func() {
		t.Lock()
		a = a + 1
		t.Unlock()
	}()

	go func() {
		t.Lock()
		a = a + 1
		t.Unlock()
	}()

	fmt.Println(a)
}
