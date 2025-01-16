/**
 * @author Manoel Ribeiro
 * @email manoel.ribeiro@unilab.edu.br
 * @create date 2025-01-16 13:54:52
 * @modify date 2025-01-16 13:55:19
 * @desc Calculate perfect number using  Euclid–Euler theorem, Mersenne primes and Miller–Rabin primality test
 */

package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	"math/rand"
)

func CheckPrimeMillerRabin(candidate *big.Int) bool {
	// Simple shortcuts.
	one := big.NewInt(1)
	two := big.NewInt(2)

	modulo := new(big.Int)
	modulo.Sub(candidate, one)

	// Write the modulo (candidate -1) number in the form
	// 2^s * d.
	s := 0
	remainder := new(big.Int)
	quotient := new(big.Int)
	quotient.Set(modulo)

	for remainder.Sign() == 0 {
		quotient.DivMod(quotient, two, remainder)
		s += 1
	}
	// The last division failed, so we must decrement `s`.
	s -= 1
	// quotient here contains the leftover which we could not divide by two,
	// and we have a 1 remaining from this last division.
	d := big.NewInt(1)
	d.Add(one, d.Mul(two, quotient))

	// Random number source for generating witnesses.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Here 10 is the precision. Every increment to this value decreases the
	// chance of a false positive by 3/4.
	for k := 0; k < 10; k++ {

		// Every witness may prove that the candidate is composite, or assert
		// nothing.
		witness := new(big.Int)
		witness.Rand(r, modulo)

		exp := new(big.Int)
		exp.Set(d)
		generated := new(big.Int)
		generated.Exp(witness, exp, candidate)

		if generated.Cmp(modulo) == 0 || generated.Cmp(one) == 0 {
			continue
		}

		for i := 1; i < s; i++ {
			generated.Exp(generated, two, candidate)

			if generated.Cmp(one) == 0 {

				return false
			}

			if generated.Cmp(modulo) == 0 {
				break
			}
		}

		if generated.Cmp(modulo) != 0 {
			// We arrived here because the `i` loop ran its course naturally
			// without meeting the `x == modulo` break.
			return false
		}
	}

	return true
}

func main() {
	var (
		p  int64
		mn *big.Int
		wg sync.WaitGroup
	)
	start := time.Now()
	if len(os.Args) == 1 {
		fmt.Println("Usage: perfectnumbers <p>")
		os.Exit(1)
	}
	pot, err := strconv.Atoi(os.Args[1])
	if err != nil || pot <= 0 {
		fmt.Println(" <p> shuld be a positive number")
		os.Exit(1)
	}
	for p = 2; p < int64(pot); p++ {
		mn = big.NewInt(0)
		mn = mn.Exp(big.NewInt(2), big.NewInt(p), nil)
		mn = mn.Sub(mn, big.NewInt(1))
		go testNumber(&wg, mn, p)
	}
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Done [%s].", elapsed)
}

func testNumber(wg *sync.WaitGroup, n *big.Int, p int64) {
	wg.Add(1)
	if CheckPrimeMillerRabin(n) {
		pot := big.NewInt(0)
		pot = pot.Exp(big.NewInt(2), big.NewInt(p-1), nil)
		pn := n.Mul(n, pot)
		fmt.Println("PN:p=", p, ",size=", len(pn.String()), pn.String())
	}
	wg.Done()
}
