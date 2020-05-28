package main

import (
	"fmt"
)

func main() {
	for {
		// fmt.Print(">>")
		// var inputInt int64
		// fmt.Scan(&inputInt)
		// fmt.Print("Steps:\t\t")
		// fmt.Println(inputInt)
		// persistenceResult := persistence(inputInt)
		// fmt.Print("Persistance:\t")
		// fmt.Println(persistenceResult)
		// fmt.Println("")
		// Git test

	}
}

// Function that gets persistance
func persistence(n int64) int8 {
	// Get next step while not 0. Keep track of persistence.
	var p int8 // p is the persistance
	p = 0
	for {
		// Get next step:
		p = p + 1
		n = doMultiply(n)
		fmt.Print("\tStep:\t")
		fmt.Println(n)
		if n < 10 {
			break
		}
	}
	return p
}

// Function that multiplies all constituent digits
func doMultiply(n int64) int64 {
	// Make empty slice
	digitsSlice := make([]int64, 0, 1000)

	var baseValue int64
	baseValue = 10

	for {
		digitsSlice = append(digitsSlice, n%baseValue)
		n = n / baseValue
		if n == 0 {
			break
		}
	}
	// Loop over function to get next step

	var returnDigit int64
	returnDigit = 1

	for _, digit := range digitsSlice {
		returnDigit = returnDigit * digit
	}
	return returnDigit
}
