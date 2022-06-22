package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// graph vertex red adjacency and black adjacency lists
	var ra, ba [][]int

	// stdin, parse, fill ra, ba
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	var accu int
	var inDigits bool
	var n int
	var one int
	var two int
	for _, b := range bytes {
		if b >= 48 && b <= 57 {
			inDigits = true
			accu = accu*10 + int(b) - 48
		} else {
			if inDigits {
				one = two
				two = accu
				accu = 0
				inDigits = false
				if n == 0 {
					n = two
					// slot[0] is never used. Index [1..n].
					ra = make([][]int, n+1)
					ba = make([][]int, n+1)
					continue
				}
			}
			switch b {
			case byte('b'):
				ba[one] = append(ba[one], two)
				ba[two] = append(ba[two], one)
				break
			case byte('r'):
				ra[one] = append(ra[one], two)
				ra[two] = append(ra[two], one)
				break
			}
		}
	}

	// Discovered-After-Red Lists
	darl := make([][]bool, n+1)

	// Do BFS from each vertex
	for s := 1; s <= n; s++ {
		// discovered
		d := make([]bool, n+1)
		// discovered after red
		dar := make([]bool, n+1)
		// fifo
		f := make([]int, 0)
		// fifo after red
		far := make([]bool, 0)
		// queue the search key
		d[s] = true
		f = append(f, s)
		far = append(far, false)
		// run until dry
		for len(f) > 0 {
			// dequeue vertex
			v := f[0]
			f = f[1:]
			// dequeue after red?
			ar := far[0]
			far = far[1:]
			for _, u := range ba[v] {
				if !d[u] {
					d[u] = true
					dar[u] = ar
					f = append(f, u)
					far = append(far, ar)
				}
			}
			for _, u := range ra[v] {
				if !d[u] {
					d[u] = true
					dar[u] = true
					f = append(f, u)
					far = append(far, true)
				}
			}
		}
		darl[s] = dar
	}

	// test the selftest
	// good, caught it: darl[3][7] = !darl[7][3]
	// as a selftest, prove the symmetry of darl
	// for i := 1; i <= n; i++ {
	// 	if darl[i][i] {
	// 		fmt.Println("darl[i][i]")
	// 	}
	// 	for j := i + 1; j <= n; j++ {
	// 		if darl[i][j] != darl[j][i] {
	// 			fmt.Println("asym")
	// 		}
	// 	}
	// }

	// Okay, I brute-forced it to here. Now what?

	N := int64(0)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if darl[i][j] {
				for k := j + 1; k <= n; k++ {
					if darl[i][k] && darl[j][k] {
						N++
					}
				}
			}
		}
	}
	fmt.Println(N)
}
