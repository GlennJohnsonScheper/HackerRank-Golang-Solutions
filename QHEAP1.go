package main

/*
Implement min-heap.

Remember what is a min-heap:
lay each generation of values in next doubling-size row of a dynamic array:
N.B. This is a 1-up numbering system
Q[0] = slot ignored
Q[1] = top element of heap
Q[2,3] = children of top...
Q[n] = final value in an n-sized heap (if n > 0)
Q[n+1] past final value

STDIN       Function
-----       --------
5           Q = 5
1 4         insert 4
1 9         insert 9
3           print minimum
2 4         delete 4
3           print minimum
*/

import (
	"fmt"
	"io/ioutil"
	//    "log"
	"os"
	//    "time"
)

func main() {
	// start := time.Now()

	// Input DIY
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	// Atoi DIY
	data := make([]int, 0)
	accu := int(0)
	negative := false
	indigits := false
	for _, b := range bytes {
		switch b {
		case 32, 13, 10, 9:
			// whitespace
			if indigits {
				if negative {
					data = append(data, -accu)
				} else {
					data = append(data, accu)
				}
				// log.Println(data[len(data)-1])
				accu = 0
				negative = false
				indigits = false
			}
			break
		case 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
			// digits
			accu = accu*10 + int(b-48)
			indigits = true
			break
		case 45:
			// hyphen
			negative = true
			break
		default:
			// log.Printf("%d\n", b)
			break
		}
	}
	// at EOF:
	if indigits {
		if negative {
			data = append(data, -accu)
		} else {
			data = append(data, accu)
		}
	}
	// process data
	state := 0
	var h Heap
	h.arr = append(h.arr, -1) // slot[0] is always ignored. Now len==1, is empty.
	for _, n := range data {
		switch state {
		case 0:
			state = 1 // ignore top line count
			break
		case 1:
			// get cmd
			switch n {
			case 1:
				state = 2 // will insert
				break
			case 2:
				state = 3 // will delete
				break
			case 3:
				// report min now
				// log.Printf("min %v\n", h.peek())
				fmt.Println(h.peek())
				break
			default:
				panic("cmd not 1,2,3")
			}
			break
		case 2:
			// insert
			// log.Printf("insert %v\n", n)
			h.insert(n)
			// log.Println(h.arr)
			state = 1 // for cmd
			break
		case 3:
			// delete
			// log.Printf("delete %v\n", n)
			h.delete(n)
			// log.Println(h.arr)
			state = 1 // for cmd
			break
		}
	}
	// log.Println(time.Since(start))
}

type Heap struct {
	arr []int
}

/*
To add a value:
add just past final value, swaps a/r up the heap.
parent = child index/2. swap if child < parent.
child = child / 2. stop when child == top.
also stop when no swap required.
*/

func (h *Heap) insert(n int) {
	h.arr = append(h.arr, n)
	// bubble up
	for ichild := len(h.arr) - 1; ichild > 1; {
		iparent := ichild / 2
		if h.arr[iparent] < h.arr[ichild] {
			break
		}
		h.arr[iparent], h.arr[ichild] = h.arr[ichild], h.arr[iparent]
		ichild = iparent
	}
}

/*
To remove top value:
move final value (at len(Q)-1) to top, swaps a/r down the heap.
child1 = parent index * 2 (if present)
child2 = parent index * 2 + 1 (if present)
for a minheap, choose the lesser child. swap if child < parent.
parent = parent * 2. stop when parent > n/2.
also stop when no swap required.

To remove any value:
find it in array, by a linear search, proceed just as for top
*/

func (h *Heap) delete(n int) {
	indexFound := -1
	for i := 1; i < len(h.arr); i++ {
		if h.arr[i] == n {
			indexFound = i
			break
		}
	}
	if indexFound == -1 {
		panic("delete not found")
	}
	// move last element into the deletion
	h.arr[indexFound] = h.arr[len(h.arr)-1]
	// shorten arr
	h.arr = h.arr[:len(h.arr)-1]
	// bubble down
	for iparent := indexFound; ; {
		ichild := iparent * 2
		// if this parent has no children
		if ichild >= len(h.arr) {
			break
		}
		// if this parent has two children, promote the minimum one
		if ichild+1 < len(h.arr) {
			if h.arr[ichild] > h.arr[ichild+1] {
				ichild++
			}
		}
		if h.arr[iparent] < h.arr[ichild] {
			break
		}
		h.arr[iparent], h.arr[ichild] = h.arr[ichild], h.arr[iparent]
		iparent = ichild
	}
}

func (h *Heap) peek() int {
	if len(h.arr) < 2 {
		panic("peek empty heap")
	}
	return h.arr[1]
}
