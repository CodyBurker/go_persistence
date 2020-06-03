// Plan:
// Create a function that decodes a morton code
// Create a function to test the persistance of morton code
// Main function to wrap all that
package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Jsewill/morton"
)

func main() {
	// Set start and end values
	startValue := uint64(0)
	endValue := uint64(100)
	// Set number of concurrent goroutines
	nThreads := 8

	getAllResults(startValue, endValue, nThreads)
}

// Function to check an entire range of numbers
func getAllResults(startValue uint64, endValue uint64, nThreads int) results {
	// Start time
	startTime := time.Now()

	// Initialize variables
	// Create morton table to decode ranges
	m := new(morton.Morton)
	m.Create(3, 512)

	// Create buffered channel to gather results
	resultsChan := make(chan results, nThreads)

	chunkSize := uint64(endValue-startValue+1) / uint64(nThreads)

	// Assign chunks to goroutines
	for i := uint64(1); i <= uint64(nThreads); i++ {
		//fmt.Println("Starting thread.")
		// For chanels 1 to (nThreads-1), assign 1/nThreads results
		// For the last channel, give them the rest of the results
		var threadStartValue uint64
		var threadendValue uint64
		if i < uint64(nThreads) {
			threadStartValue = startValue + chunkSize*(i-1)
			threadendValue = threadStartValue + chunkSize - 1
		} else {
			threadStartValue = startValue + chunkSize*(i-1)
			threadendValue = endValue
		}
		go getResults(resultsChan, threadStartValue, threadendValue, *m)
	}
	// Gather results as threads finish
	finalResults := results{0, 0, *big.NewInt(0)}
	for i := 0; i < nThreads; i++ {
		threadResults := <-resultsChan
		if threadResults.maxPersistenceValue > finalResults.maxPersistenceValue {
			threadResults.maxPersistenceValue = finalResults.maxPersistenceValue
			threadResults.maxPersistenceNumber = finalResults.maxPersistenceNumber
		}
		finalResults.totalPeristence.Add(&threadResults.totalPeristence, &finalResults.totalPeristence)
	}
	elapsed := time.Since(startTime)
	fmt.Printf("Threads:%d\tTime:%s", nThreads, elapsed)
	fmt.Printf("Results:\n\tMaxPersistenceValue:\t%d\n\tTotalPersistence:\t%s", finalResults.maxPersistenceValue, finalResults.totalPeristence.String())
	return finalResults
}

// Function to check a chunk of numbers
func getResults(resultsChan chan results, startValue uint64, endValue uint64, decoder morton.Morton) {

	threadResults := results{0, 0, *big.NewInt(0)}
	totalPeristence := big.NewInt(0)
	// Loop over assigned chunk of numbers, and gather statistics
	for i := startValue; i <= endValue; i++ {
		startValue := decoder.Decode(i)
		x := startValue[0]
		y := startValue[1]
		z := startValue[2]
		loopNumber := convertFactors(x, y, z)
		loopResults := getPersistence(loopNumber) + 1
		if loopResults > threadResults.maxPersistenceValue {
			threadResults.maxPersistenceValue = loopResults
			threadResults.maxPersistenceNumber = i
		}
		totalPeristence.Add(totalPeristence, big.NewInt(int64(loopResults)))
	}
	threadResults.totalPeristence = *totalPeristence
	fmt.Println(threadResults)
	//fmt.Printf("Thread done evaluating:\t%d\t%d\n\t(%s\t%d\t%d)\n", startValue, endValue, threadResults.totalPeristence.String(), threadResults.maxPersistenceNumber, threadResults.maxPersistenceValue)
	resultsChan <- threadResults
	return
}

// A struct to pass results back to the main loop
type results struct {
	maxPersistenceValue  int
	maxPersistenceNumber uint64
	totalPeristence      big.Int
}

//  f(x,y,z) = (2 ^ x) * (3 ^ y) * (7 ^ z)
func convertFactors(x, y, z uint32) (result *big.Int) {
	// result=1
	result = big.NewInt(1)
	// result= result * 2 ^ x
	result.Mul(result, big.NewInt(1).Exp(big.NewInt(2), big.NewInt(int64(x)), nil))

	// result=result * 3 ^ y
	result.Mul(result, big.NewInt(1).Exp(big.NewInt(3), big.NewInt(int64(y)), nil))

	// result=result * 7 ^ z
	result.Mul(result, big.NewInt(1).Exp(big.NewInt(7), big.NewInt(int64(z)), nil))
	return
}

// Get multiplicitive persitence of a number
func multiplyDigits(inputNum *big.Int) *big.Int {

	// The product of the digits to return
	result := big.NewInt(1)
	digit := big.NewInt(0)
	// We are working in base 10
	for {
		base := big.NewInt(10)
		// Get digit, remainder
		inputNum, digit = inputNum.DivMod(inputNum, base, base)
		// fmt.Printf("Intput:\t%s\nDigit:\t%s\nResult:\t%s\n", inputNum.String(), digit.String(), result.String())
		// fmt.Print("Test:")
		// fmt.Println(digit.Cmp(big.NewInt(0)))
		// If digit is <= 0 but inputNum >0, there is a zero digit
		if digit.Cmp(big.NewInt(0)) <= 0 {
			if inputNum.Cmp(big.NewInt(0)) > 0 {
				result := big.NewInt(0)
				return result
			}
			// If digit <= 0 and result inputNum <= 0, we are at the end of the number - return results
			return result

		}
		// Otherwise, multiply result by digit
		result.Mul(result, digit)

	}

}

func getPersistence(inputNum *big.Int) (returnInt int) {
	returnInt = 0
	base := big.NewInt(10)
	for {
		//fmt.Println(inputNum)
		if inputNum.Cmp(base) < 0 {
			return
		}
		inputNum = multiplyDigits(inputNum)
		//		fmt.Println(inputNum)
		returnInt = returnInt + 1
	}

}
