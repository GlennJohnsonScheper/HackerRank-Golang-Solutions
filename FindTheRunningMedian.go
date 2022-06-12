// Hackerrank 2022-06-08 task (GO) Find the running median.txt

// https://www.hackerrank.com/challenges/find-the-running-median/problem

// stdin = 1 positive count < 10^5, N integers, non-negative < 10^5.
// stdout = N medians = center element or average of two, to 1 decimal place.
// written from scratch in Go, from memory with a little help from google.
// started 4:54 pm.
// finished 10:11 pm.
// Plan: Keep a min-heap and a max-heap on either side of the middle.

package main

import (
	"fmt"
	"io/ioutil"
	// "log"
	"os"
)

type Heap struct {
	arr  []int
	test func(a, b int) bool // minheap=less, or maxheap=more
}

func (h *Heap) insert(n int) {
	// append new item past last item
	h.arr = append(h.arr, n)
	// rebalance up
	for child := len(h.arr) - 1; ; {
		if child < 2 {
			break
		}
		parent := child / 2
		if h.test(h.arr[child], h.arr[parent]) {
			h.arr[child], h.arr[parent] = h.arr[parent], h.arr[child]
		} else {
			break
		}
		child = parent
	}
	// log.Println("inserted", n, "so", len(h.arr), h.arr)
}

func (h *Heap) pop() {
	logDetail := h.arr[1]
	// copy last item over top item
	h.arr[1] = h.arr[len(h.arr)-1]
	h.arr = h.arr[:len(h.arr)-1]
	// rebalance down
	last := len(h.arr) - 1
	for parent := 1; ; {
		child := parent * 2
		if child > last {
			// There is no child
			break
		}
		if child < last {
			// There are two children
			// Which child?
			// suppose it is for minheap,
			// and test func returns a<b.
			// for test==true, a is smaller, and goes up
			// for test==false, b is smaller, and goes up
			if !h.test(h.arr[child], h.arr[child+1]) {
				child++ // work with the b child
			}
		}
		if h.test(h.arr[child], h.arr[parent]) {
			h.arr[child], h.arr[parent] = h.arr[parent], h.arr[child]
		} else {
			break
		}
		parent = child
	}
	_ = logDetail
	// log.Println("popped", logDetail, "so", len(h.arr), h.arr)
}

type Passer struct {
	minheap *Heap
	maxheap *Heap
}

func main() {
	// log.Println("HelloWorld")
	var args Passer

	// stdin: DIY
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic("Stdin ReadAll Error")
	}

	// atoi: DIY
	var gotCount bool
	var inDigits bool
	var accu int
	for _, b := range bytes {
		// log.Println(b)
		switch b {
		case 48, 49, 50, 51, 52, 53, 54, 55, 56, 57: // digits
			accu = accu*10 + int(b-48)
			inDigits = true
			break
		case 32, 10, 13, 9: // whitespace
			if inDigits {
				if gotCount {
					doNumber(&args, accu)
				} else {
					gotCount = true
					// Process the top line of pos. count:
					// Create one minheap and one maxheap,
					// both with capacity for all numbers.
					// Note: heap[0] stays forever unused.

					args.minheap = &Heap{arr: make([]int, 1, accu+1), test: func(a, b int) bool { return a < b }}
					args.maxheap = &Heap{arr: make([]int, 1, accu+1), test: func(a, b int) bool { return a > b }}
				}
				accu = 0
			}
			inDigits = false
			break
		}
	}
	// at EOF w/o newline
	if inDigits {
		doNumber(&args, accu)
	}
}

func doNumber(args *Passer, n int) {
	// log.Println("doing", n)
	// in case of an odd count, let the minheap be larger.
	// otherwise, transfer a/r to maintain an equal count.
	// A count of ONE is the new ZERO as indexing is 1-up.
	if len(args.minheap.arr) == 1 {
		// log.Println("initial insert to minheap")
		args.minheap.insert(n)
	} else {
		// Henceforth, ZERO is moot,
		// I need to peek() minheap,
		// else appending to either.

		// N.B. minheap is holding the LARGE values
		// N.B. maxheap is holding the SMALL values
		if n > args.minheap.arr[1] {
			// log.Println("insert to minheap")
			args.minheap.insert(n)
		} else {
			// log.Println("insert to maxheap")
			args.maxheap.insert(n)
		}
		// Rebalance.
		both := len(args.minheap.arr) + len(args.maxheap.arr)
		if both/2 != len(args.maxheap.arr) {
			// Just an If, not a While:
			// Only ever move one item.
			if both/2 < len(args.maxheap.arr) {
				// move one item from maxheap to minheap
				// log.Println("move one item from maxheap to minheap")
				args.minheap.insert(args.maxheap.arr[1])
				args.maxheap.pop()
			} else {
				// move one item from minheap to maxheap
				// log.Println("move one item from minheap to maxheap")
				args.maxheap.insert(args.minheap.arr[1])
				args.minheap.pop()
			}
		}
	}
	// N.B. items are not sorted within every one doubling layer:
	// log.Println("minheap", len(args.minheap.arr), args.minheap.arr)
	// log.Println("maxheap", len(args.maxheap.arr), args.maxheap.arr)

	var ans int
	if len(args.minheap.arr) > len(args.maxheap.arr) {
		ans = 10 * args.minheap.arr[1]
	} else {
		ans = 5 * (args.minheap.arr[1] + args.maxheap.arr[1])
	}
	fmt.Printf("%v.%v\n", ans/10, ans%10)

}
