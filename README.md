# go_persistence
## Background
Multiplicative persistence is the process of taking any natural number, and then multiplying its digits together. Afterwards one can take the resulting number, and repeat the process until only a single digit remains. The number of times this process can done before reaching a single digit (or zero) is called its multiplicative persistence. [Read more on Wolfram Alpha.](https://mathworld.wolfram.com/MultiplicativePersistence.html)

It is believed that the largest multiplicative persistence of any number is 11: [link](https://listserv.nodak.edu/cgi-bin/wa.exe?A2=ind0107&L=NMBRTHRY&P=R1036&I=-3). It is worth noting that not every number need be checked sequentially. Some numbers can be skipped: for example, any number with a 0 will have a persistence of 0. Any number with a 2 and a 5 will have a persistence of 2 (the resulting number will end in a 0). Any number with a 1 in it wil have a smaller number that has the same persistence. 

Carmody ([link](https://listserv.nodak.edu/cgi-bin/wa.exe?A2=ind0107&L=NMBRTHRY&P=R1036&I=-3)) points out that one need only check numbers that consist of factors of 2,3,7 or 3,5,7. This program provides the function `getAllResults` and various supporting functions that check various combinations of 2,3,7 sequentially and concurrently. The intended use of this function is to provide a fast and sequential way to check a large set of numbers. In order to sequence combinations of these factors (2,3,7) this code utilizes Morton encoding ([Wikipedia](https://en.wikipedia.org/wiki/Z-order_curve)) to map a single dimensional space (in this case a `uint64`) to three dimensional space (with each dimension corresponding to a factor of 2,3,7). The Jsewell Morton library is used for this ([Github link](https://github.com/Jsewill/morton)). This provides a method to check the 'first' combination of factors, then the 'second', etc. and theoretically exhaust the space. 

## Use
To use the main function provided by this code,`getAllResults`, provide a start value and end value for the ranges of the morton code you want to check. Then provide the number of threads (read: gofunctions) you wish to utilize. Preliminary testing suggests to set this to the number of hyperthreads on your CPU (e.g. 8). The result will be a `struct` of the type `results`, which has the properties of the largest multiplicative persistence found in the range you provided, the number that provided that persistence, and finally the sum of the multiplicative persistence of all of the numbers.

## Hilights
Some hilights of this code include:
* Concurrency (using go functions) 
* Custom data structures (`struct`s)
* Using Go's `math/big` library to handle unconventionally large numbers.
