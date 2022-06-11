// MinimumAverageWaitingTime.go

/*
Hackerrank 2022-06-10 task (GO) Minimum Average Waiting Time

Starting at 2:45 pm
Finished at 9:00 pm

What could they mean by:
The i-th customer is not the customer arriving at the i-th arrival time.
I take that to mean that the arrival times are not necessarily in order.
Therefore, I must pre-input all pairs and pre-sort them by arrival time.

Strategy:
At any moment when a pizza is done,
look at all then available choices,
(because cook has no foreknowledge)
and choose the shortest pizza time.

I suppose if any pizza ends at time n,
and an arrival happens at same time n,
that arrival is eligible to be chosen.

Keep a minheap of all available choices.
Add to minheap all from array <= time n.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type Pair struct {
	arrival int64
	cooking int64
}

type MinHeap struct {
	arr []Pair // sorted by cooking, not arrival!
}

func (mh *MinHeap) insert(p Pair) {
	mh.arr = append(mh.arr, p)
	for child := len(mh.arr) - 1; child > 1; {
		parent := child / 2
		if mh.arr[child].cooking < mh.arr[parent].cooking {
			mh.arr[child], mh.arr[parent] = mh.arr[parent], mh.arr[child]
		}
		child = parent
	}
}

func (mh *MinHeap) peek() Pair {
	return mh.arr[1]
}

func (mh *MinHeap) pop() {
	mh.arr[1] = mh.arr[len(mh.arr)-1]
	mh.arr = mh.arr[:len(mh.arr)-1]
	for parent := 1; ; {
		child := parent * 2
		if child > len(mh.arr)-1 {
			return
		}
		if child < len(mh.arr)-1 {
			//Which child?
			if mh.arr[child+1].cooking < mh.arr[child].cooking {
				child++
			}
		}
		if mh.arr[child].cooking < mh.arr[parent].cooking {
			mh.arr[child], mh.arr[parent] = mh.arr[parent], mh.arr[child]
		}
		parent = child
	}
}

func main() {
	// Stdin, Atoi - DIY:
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var haveCount bool
	var accu int64
	var inDigits bool
	var arrival int64
	var haveArrival bool
	buf := make([]Pair, 0)

	for _, b := range bytes {
		switch b {
		case 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
			inDigits = true
			accu = accu*10 + int64(b) - 48
			break
		case 32, 10, 13:
			if inDigits {
				inDigits = false
				if haveCount == false {
					haveCount = true
					accu = 0
					inDigits = false
					continue // skip 1st integer
				}
				if haveArrival {
					haveArrival = false
					buf = append(buf, Pair{arrival: arrival, cooking: accu})
				} else {
					haveArrival = true
					arrival = accu
				}
				accu = 0
			}
			break
		}
	}
	if inDigits {
		buf = append(buf, Pair{arrival: arrival, cooking: accu})
	}
	sort.Slice(buf, func(i, j int) bool { return buf[i].arrival < buf[j].arrival })

	var clock int64
	var sumWait int64
	var qtyWait int64
	mh := MinHeap{arr: make([]Pair, 1)} // Slot 0 is never used. Len()==1 means empty.

	// run once through buffer
	for i := 0; ; {
		// phase 1: heap up while arrival <= clock
		for i < len(buf) {
			if buf[i].arrival <= clock || len(mh.arr) == 1 {
				mh.insert(buf[i])
				i++
			} else {
				break
			}
		}
		// phase 2: serve one atop heap
		// Why not more? Go around loop.
		if len(mh.arr) > 1 {
			p := mh.peek()
			if clock < p.arrival {
				clock = p.arrival
				sumWait += p.cooking // no waiting to start cooking
			} else {
				sumWait += clock - p.arrival + p.cooking // waiting + cooking
			}
			clock += p.cooking
			qtyWait++
			mh.pop()
		}
		// When is loop done?
		if i == len(buf) && len(mh.arr) == 1 {
			break
		}
	}
	// slight faux pas here: HR template obeys os.Getenv("OUTPUT_PATH")
	fmt.Printf("%v\n", sumWait/qtyWait)
}
