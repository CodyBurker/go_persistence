package main

import (
	"fmt"
)

func main() {
	persistenceTest()
	// for {
	// fmt.Print(">>")
	// var inputInt int64
	// fmt.Scan(&inputInt)
	// fmt.Print("Steps:\t\t")
	// fmt.Println(inputInt)
	// persistenceResult := persistence(inputInt)
	// fmt.Print("Persistance:\t")
	// fmt.Println(persistenceResult)
	// fmt.Println("")
	// }
}

// TODO: Write a test to make sure that my changes don't mess with it
func persistenceTest() {
	// Use known results from here: https://mathworld.wolfram.com/MultiplicativePersistence.html
	testOutputs := [11]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	testInputs := [11]int{10, 25, 39, 77, 679, 6788, 68889, 2677889, 26888999, 3778888999, 277777788888899}
	for index, inputDigit := range testInputs {
		if persistence(int64(inputDigit)) != int8(testOutputs[index]) {
			fmt.Println("Error:")
			fmt.Print("Input:\t")
			fmt.Println(inputDigit)
			fmt.Print("Got:\t")
			fmt.Println(persistence(int64(inputDigit)))
			fmt.Print("Expected:\t")
			fmt.Println(testOutputs[index])
		}
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
		// fmt.Print("\tStep:\t")
		// fmt.Println(n)
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
		newDigit := n % baseValue
		digitsSlice = append(digitsSlice, newDigit)
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
