package main

// Need to implement

// https://marcin-chwedczuk.github.io/iterative-algorithm-for-drawing-hilbert-curve
// http://www.math.uwaterloo.ca/~wgilbert/Research/HilbertCurve/HilbertCurve.html
// https://listserv.nodak.edu/cgi-bin/wa.exe?A2=ind0107&L=NMBRTHRY&P=R1036&I=-3
// https://en.wikipedia.org/wiki/Z-order_curve

// Plan: Use this to decode an incrementing binary
//https://github.com/Jsewill/morton

import (
	"fmt"
	"math"
	"time"
)

func main() {
	persistenceTest()

	for {
		fmt.Print(">>")
		var inputInt uint64 = 0
		var maxPersistance int8 = 0
		start := time.Now()
		cycle := 0
		max := 500000
		for {
			inputInt = nextNumber(inputInt)
			result := persistence(inputInt)
			if result > maxPersistance {
				fmt.Printf("%d\t\t\t%d\n", inputInt, result)
				maxPersistance = result
			}
			if result > 11 {
				fmt.Println("WHAT??")
				fmt.Println(inputInt)
				fmt.Println(result)
				break
			}
			cycle++
			if cycle >= max {
				cycle = 0
				end := time.Now()
				fmt.Printf("Avg Time: %d\tinput:%d\tresult:%d\n", int(end.Sub(start))/max, inputInt, maxPersistance)
			}

		}
	}
}

// Function that generates the next candidate number
func nextNumber(lastNumber uint64) uint64 {
	// Check if numbers are in ascending order
	// If not, then increment the next digit
	var candidate uint64
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
func getNumber(digitsSlice []uint64) uint64 {
	var returnValue uint64 = 0
	for index, thisDigit := range digitsSlice {
		returnValue = returnValue + uint64(math.Pow10(index))*uint64(thisDigit)
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
		if persistence(uint64(inputDigit)) != int8(testOutputs[index]) {
			fmt.Println("Error:")
			fmt.Print("Input:\t")
			fmt.Println(inputDigit)
			fmt.Print("Got:\t")
			fmt.Println(persistence(uint64(inputDigit)))
			fmt.Print("Expected:\t")
			fmt.Println(testOutputs[index])
		}
	}
}

// Function that gets persistance
func persistence(n uint64) int8 {
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

func getDigits(n uint64) []uint64 {
	// Initialize empty slice
	digitsSlice := make([]uint64, 0, 1000)
	var baseValue uint64
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
func doMultiply(n uint64) uint64 {
	// Make empty slice
	digitsSlice := getDigits(n)
	// Loop over function to get next step

	var returnDigit uint64
	returnDigit = 1

	for _, digit := range digitsSlice {
		returnDigit = returnDigit * digit
	}
	return returnDigit
}
