package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	// Number of threads
	numberOfThreads := 1000
	resultsSlice := make(chan results, numberOfThreads)
	var wg sync.WaitGroup
	wg.Add(numberOfThreads)
	startValue := 1
	endValue := 100000000

	length := (1 + endValue - startValue) / numberOfThreads

	// Split up data, and assign to threads
	for i := 1; i <= numberOfThreads; i++ {
		// For chanels 1 to (nThreads-1), assign 1/nThreads results
		// For the last channel, give them the rest of the results
		var threadStartValue int
		var threadendValue int

		if i < numberOfThreads {
			threadStartValue = startValue + length*(i-1)
			threadendValue = threadStartValue + length - 1
		} else {
			threadStartValue = startValue + length*(i-1)
			threadendValue = endValue
		}

		go getResults(resultsSlice, threadStartValue, threadendValue, &wg)
	}
	// Gather results as they come in
	maxPersistenceValue := 0
	maxPeristence := 0
	totalPeristence := 0
	// Read from buffered channel when ready.
	for i := 1; i <= numberOfThreads; i++ {
		newResult := <-resultsSlice
		if newResult.maxPeristence > maxPeristence {
			maxPersistenceValue = newResult.maxPersistenceValue
			maxPeristence = newResult.maxPeristence
		}
		totalPeristence = totalPeristence + newResult.totalPeristence
	}
	fmt.Printf("maxP:\t%d\nmaxPV:\t%d\nTotalP:\t%d", maxPeristence, maxPersistenceValue, totalPeristence)
}

type results struct {
	maxPersistenceValue int
	maxPeristence       int
	totalPeristence     int
}

func getResults(inputResults chan results, startValue int, endvalue int, wg *sync.WaitGroup) {
	fmt.Printf("\tThread from %d to %d started.\n", startValue, endvalue)
	//time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
	theseResults := results{rand.Int(), rand.Int(), endvalue - startValue + 1}
	inputResults <- theseResults
	fmt.Printf("\t\t\tThread %d done.\n", startValue)
}
