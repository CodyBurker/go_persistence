// Plan:
// Create a function that decodes a morton code
// Create a function to test the persistance of morton code
// Main function to wrap all that
package main

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/Jsewill/morton"
)

func main() {

	// Start time
	startTime := time.Now()

	// Initialize variables
	// Create morton table to decode ranges
	m := new(morton.Morton)
	m.Create(3, 512)
	// Set number of concurrent goroutines
	nThreads := 8
	// Create buffered channel to gather results
	resultsChan := make(chan results, nThreads)
	// Set start and end values
	startValue := uint64(0)
	endValue := uint64(10000)
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

// Can't figure out how to properly install this:

// Below code belongs to Jsewill at https://github.com/Jsewill/morton
// I have no idea how to install the package properly so for now I am copying the code into my file
type Table struct {
	Index  uint8
	Length uint32
	Encode []Bit
}

// Sortable Table slice type to satisfy the sort package interface
type ByTable []Table

func (t ByTable) Len() int {
	return len(t)
}

func (t ByTable) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByTable) Less(i, j int) bool {
	return t[i].Index < t[j].Index
}

type Bit struct {
	Index uint32
	Value uint64
}

// Sortable Table slice type to satisfy the sort package interface
type ByBit []Bit

func (b ByBit) Len() int {
	return len(b)
}

func (b ByBit) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByBit) Less(i, j int) bool {
	return b[i].Index < b[j].Index
}

type Morton struct {
	Dimensions uint8
	Tables     []Table
	Magic      []uint64
}

// Convenience function
func New(dimensions uint8, size uint32) *Morton {
	m := new(Morton)
	m.Create(dimensions, size)
	return m
}

func (m *Morton) Create(dimensions uint8, size uint32) {
	done := make(chan struct{})
	mch := make(chan []uint64)
	go func() {
		m.CreateTables(dimensions, size)
		done <- struct{}{}
	}()
	go func() {
		mch <- MakeMagic(dimensions)
	}()
	m.Magic = <-mch
	close(mch)
	<-done
	close(done)
}

func (m *Morton) CreateTables(dimensions uint8, length uint32) {
	ch := make(chan Table)

	m.Dimensions = dimensions
	for i := uint8(0); i < dimensions; i++ {
		go func(i uint8) {
			ch <- CreateTable(i, dimensions, length)
		}(i)
	}
	for i := uint8(0); i < dimensions; i++ {
		t := <-ch
		m.Tables = append(m.Tables, t)
	}
	close(ch)

	sort.Sort(ByTable(m.Tables))
}

func MakeMagic(dimensions uint8) []uint64 {
	// Generate nth and ith bits variables
	d := uint64(dimensions)
	limit := 64/d + 1
	nth := []uint64{0, 0, 0, 0, 0}
	for i := uint64(0); i < 64; i++ {
		if i < limit {
			nth[0] |= 1 << (i * (d))
		}

		nth[1] |= 3 << (i * (d << 1))
		nth[2] |= 0xf << (i * (d << 2))
		nth[3] |= 0xff << (i * (d << 3))
		nth[4] |= 0xffff << (i * (d << 4))
	}

	return nth
}

func (m *Morton) Decode(code uint64) (result []uint32) {
	if m.Dimensions == 0 {
		return
	}

	d := uint64(m.Dimensions)
	r := make([]uint64, d)

	// Process each dimension
	for i := uint64(0); i < d; i++ {
		r[i] = (code >> i) & m.Magic[0]

		r[i] = (r[i] ^ (r[i] >> (1 << (d - 2)))) & m.Magic[1]
		r[i] = (r[i] ^ (r[i] >> (2 << (d - 2)))) & m.Magic[2]
		r[i] = (r[i] ^ (r[i] >> (4 << (d - 2)))) & m.Magic[3]
		r[i] = (r[i] ^ (r[i] >> (8 << (d - 2)))) & m.Magic[4]

		result = append(result, uint32(r[i]))
	}

	return
}

func (m *Morton) Encode(vector []uint32) (result uint64, err error) {
	length := len(m.Tables)
	if length == 0 {
		err = errors.New("No lookup tables.  Please generate them via CreateTables().")
		return
	}

	if len(vector) > length {
		err = errors.New("Input vector slice length exceeds the number of lookup tables.  Please regenerate them via CreateTables()")
		return
	}

	//sort.Sort(sort.Reverse(ByUint32Index(vector)))

	for k, v := range vector {
		if v > uint32(len(m.Tables[k].Encode)-1) {
			err = errors.New(fmt.Sprint("Input vector component, ", k, " length exceeds the corresponding lookup table's size.  Please regenerate them via CreateTables() and specify the appropriate table length"))
			return
		}

		result |= m.Tables[k].Encode[v].Value
	}

	return
}

func CreateTable(index, dimensions uint8, length uint32) Table {
	t := Table{Index: index, Length: length}
	bch := make(chan Bit)

	// Build interleave queue
	for i := uint32(0); i < length; i++ {
		go func(i uint32) {
			bch <- InterleaveBits(i, uint32(index), uint32(dimensions-1))
		}(i)
	}
	// Pull from interleave queue
	for i := uint32(0); i < length; i++ {
		ib := <-bch
		t.Encode = append(t.Encode, ib)
	}
	close(bch)

	sort.Sort(ByBit(t.Encode))
	return t
}

// Interleave bits of a uint32.
func InterleaveBits(value, offset, spread uint32) Bit {
	ib := Bit{value, 0}

	// Determine the minimum number of single shifts required. There's likely a better, and more efficient, way to do this.
	n := value
	limit := uint64(0)
	for i := uint32(0); n != 0; i++ {
		n = n >> 1
		limit++
	}

	// Offset value for interleaving and reconcile types
	v, o, s := uint64(value), uint64(offset), uint64(spread)
	for i := uint64(0); i < limit; i++ {
		// Interleave bits, bit by bit.
		ib.Value |= (v & (1 << i)) << (i * s)
	}
	ib.Value = ib.Value << o

	return ib
}
