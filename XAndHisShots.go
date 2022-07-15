// XAndHisShots.go

// https://www.hackerrank.com/challenges/x-and-his-shots/problem

// There are N first number ranges from Ai to Bi.
// There are M other number ranges from Ci to Di.
// Count all the overlaps of any A-B and any C-D.
// M,N in 1..10^5
// A-D in 1..10^8

// Insight 1: Way easier to formulate NOT rule:
// if A > D || C > B then these do NOT overlap.

// Insight 2: Sort all As, all Bs, all Cs, all Ds;
// Run two cursors over A and D, and over C and B;
// For any item in A, all lower Ds do not overlap.
// For any item in C, all lower Bs do not overlap.
// Subtract non-overlap count from possible (M*N).

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func sortMapToKeysValues(m map[int]int) ([]int, []int) {
	k := make([]int, 0, len(m))
	v := make([]int, 0, len(m))
	for key, _ := range m {
		k = append(k, key)
	}
	sort.Ints(k)
	for _, key := range k {
		v = append(v, m[key])
	}
	return k, v
}

func main() {
	// server logged an error to accept my submission.
	// Maybe faux pas of using stdout directly did it.
	// So reuse this HackerRank template Getenv idiom.
	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	if err != nil {
		panic(err)
	}
	defer stdout.Close()
	N := 0
	M := 0
	a := make(map[int]int)
	b := make(map[int]int)
	c := make(map[int]int)
	d := make(map[int]int)
	// closure over main variables
	solve := func() {
		ka, va := sortMapToKeysValues(a)
		kb, vb := sortMapToKeysValues(b)
		kc, vc := sortMapToKeysValues(c)
		kd, vd := sortMapToKeysValues(d)
		overlaps := int64(M) * int64(N)
		// Count any A > any D
		{
			runts := int64(0) // Running count of small Ds
			ia := 0
			id := 0
			// bi=phase loop while As remain
			for ia < len(ka) {
				testa := ka[ia]
				// For every one A loop, do zero to many D:
				for id < len(kd) {
					testd := kd[id]
					if testd < testa {
						runts += int64(vd[id])
						id++
					} else {
						break
					}
				}
				overlaps -= runts * int64(va[ia])
				ia++
			}
		}
		// Count any C > any B
		{
			runts := int64(0) // Running count of small Bs
			ic := 0
			ib := 0
			// bi=phase loop while Cs remain
			for ic < len(kc) {
				testc := kc[ic]
				// For every one C loop, do zero to many B:
				for ib < len(kb) {
					testb := kb[ib]
					if testb < testc {
						runts += int64(vb[ib])
						ib++
					} else {
						break
					}
				}
				overlaps -= runts * int64(vc[ic])
				ic++
			}
		}
		// fmt.Println(overlaps) // answer
		fmt.Fprintf(stdout, "%v\n", overlaps)
	}

	// diy input, Atoi
	accu := 0
	inDigit := false
	state := 0
	lo := 0
	hi := 0
	pairs := 0
	// closure over main variables
	useAccu := func() {
		switch state {
		case 0:
			N = accu
			state = 1
			break
		case 1:
			M = accu
			state = 2
			break
		case 2:
			lo = accu
			state = 3
			break
		case 3:
			hi = accu
			if lo > hi {
				lo, hi = hi, lo
			}
			if _, ok := a[lo]; ok {
				a[lo]++
			} else {
				a[lo] = 1
			}
			if _, ok := b[hi]; ok {
				b[hi]++
			} else {
				b[hi] = 1
			}
			state = 2
			pairs++
			if pairs == N {
				state = 4
				pairs = 0
			}
			break
		case 4:
			lo = accu
			state = 5
			break
		case 5:
			hi = accu
			if lo > hi {
				lo, hi = hi, lo
			}
			if _, ok := c[lo]; ok {
				c[lo]++
			} else {
				c[lo] = 1
			}
			if _, ok := d[hi]; ok {
				d[hi]++
			} else {
				d[hi] = 1
			}
			state = 4
			pairs++
			if pairs == M {
				solve()
				state = 6
			}
			break
		default:
			panic("state")
			break
		}
	}

	// drive the state machines
	sb, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	for _, b := range sb {
		if b < byte('0') {
			if inDigit {
				useAccu()
				inDigit = false
				accu = 0
			}
		} else {
			accu = accu*10 + int(b-byte('0'))
			inDigit = true
		}
	}
	if inDigit {
		useAccu() // at EOF
	}
}
