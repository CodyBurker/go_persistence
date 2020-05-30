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
	"sync"

	"github.com/Jsewill/morton"
)

func main() {
	// Morton decoding
	m := new(morton.Morton)
	m.Create(3, 512)
	startValue1, _ := m.Encode([]uint32{12, 7, 2})
	startValue := startValue1 - 100
	endValue := startValue1 + 500

	for i := startValue; i < endValue; i++ {

	}
	var maxPersistence chan int = make(chan int, 0)
	var maxPeristenceValue chan uint64 = make(chan uint64, startValue)
	var totalPeristence chan uint64 = make(chan uint64, 0)

	// Add a waitgroup
	var wg sync.WaitGroup
	wg.Add(int(1 + endValue - startValue))

	for i := startValue; i < endValue; i++ {
		go checkNumber(m, i, maxPersistence, maxPeristenceValue, totalPeristence)
	}
	wg.Wait()
	fmt.Println(<-maxPersistence)
	fmt.Println(<-maxPeristenceValue)
	fmt.Println(<-totalPeristence)
}

//// Function that given start and end values, returns statistics, max, etc.
// func checkValues(startValue, endvalue uint64) (maxPersistenceValue uint64, maxPersistence int, avgP float64) {
// 	// Create tables, magic bits for decoding morton to generate coordinates.
// 	m := new(morton.Morton)
// 	m.Create(3, 512)

// 	// Variable to hold largest persistence across all checked values
// 	maxPersistence = 0
// 	// Variable to hold total persistence across all values (to use for average persistence)
// 	totalPersistence := uint64(0)
// 	// Variable to hold mortoncode value of vector with largest peristence found
// 	maxPersistenceValue = startValue
// 	// No idea why this is needed.
// 	endValue := endvalue
// 	// Loop over morton code values, checking each.
// 	for i := startValue; i <= endValue; i++ {
// 		// Decode values, save to variables for legibility
// 		exponents := m.Decode(i)
// 		x := exponents[0]
// 		y := exponents[1]
// 		z := exponents[2]
// 		// Turn the above exponent into a number
// 		startValue := convertFactors(x, y, z)
// 		// Get peristence of these numbers
// 		persistence := getPersistence(startValue) + 1 // convertFactors finds the first level of persistence
// 		// Keep track of total persistence
// 		totalPersistence = totalPersistence + uint64(persistence)
// 		if persistence >= maxPersistence {
// 			maxPersistence = persistence
// 			maxPersistenceValue = i
// 		}
// 		// // Print to command line
// 		// fmt.Printf("p(%d,%d,%d)=\t%d\n", x, y, z, persistence)
// 	}
// 	// Average persistence across all checked values
// 	avgP = float64(totalPersistence) / float64(endValue-startValue)
// 	// // Print totals to command line:
// 	// fmt.Printf("\tMaxValue:\t%d\n\tMaxP    :\t%d\n\tAvgP    :\t%f",
// 	// 	maxPersistenceValue,
// 	// 	maxPersistence,
// 	// 	avgP)
// 	return
// }

// Concurrency function to check a single number, and return results to the channel
func checkNumber(m *morton.Morton, number uint64, maxPersistence chan int, maxPeristenceValue chan uint64, totalPeristence chan uint64) {
	exponents := m.Decode(number)
	x := exponents[0]
	y := exponents[1]
	z := exponents[2]
	product := convertFactors(x, y, z)
	persistence := getPersistence(product)
	if persistence >= <-maxPersistence {
		maxPersistence <- persistence
		maxPeristenceValue <- number
	}
	totalPeristence <- uint64(persistence) + <-totalPeristence
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
