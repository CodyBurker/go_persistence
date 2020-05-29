package main

import (
	"fmt"
	"math"
)

func main() {
	persistenceTest()

	for {
		fmt.Print(">>")
		var inputInt int64 = 0
		var maxPersistance int8 = 0
		for {
			inputInt = nextNumber(inputInt)
			result := persistence(inputInt)
			if result > maxPersistance {
				fmt.Printf("%d\t\t\t%d\n", inputInt, result)
				maxPersistance = result
			}
		}
	}
}

// Function that generates the next candidate number
func nextNumber(lastNumber int64) int64 {
	// Check if numbers are in ascending order
	// If not, then increment the next digit
	var candidate int64
	candidate = lastNumber + 1
	digitsSlice := getDigits(candidate)

	// Replace all 1's and 0's with 2's
	for index, thisDigit := range digitsSlice {
		if (thisDigit == 0) || (thisDigit == 1) {
			digitsSlice[index] = 2
		}
	}

	// If any digit is not ascending, increase it to the previous digit
	for i := len(digitsSlice) - 1; i >= 1; i-- {
		thisDigit := digitsSlice[i]
		nextDigit := digitsSlice[i-1]
		if thisDigit > nextDigit {
			digitsSlice[i-1] = thisDigit
		}
	}
	return getNumber(digitsSlice)
}

// Get number from slice of digits, with lowest first
func getNumber(digitsSlice []int64) int64 {
	var returnValue int64 = 0
	for index, thisDigit := range digitsSlice {
		returnValue = returnValue + int64(math.Pow10(index))*int64(thisDigit)
		// Debugging:
		//fmt.Printf("\tPower:\t%d\tReturnValue:\t%d\tDigit:\t%d\n", power, returnValue, thisDigit)
	}
	return returnValue
}

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

func getDigits(n int64) []int64 {
	// Initialize empty slice
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
	return digitsSlice
}

// Function that multiplies all constituent digits
func doMultiply(n int64) int64 {
	// Make empty slice
	digitsSlice := getDigits(n)
	// Loop over function to get next step

	var returnDigit int64
	returnDigit = 1

	for _, digit := range digitsSlice {
		returnDigit = returnDigit * digit
	}
	return returnDigit
}
