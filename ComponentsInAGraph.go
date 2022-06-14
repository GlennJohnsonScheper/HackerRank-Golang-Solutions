// Hackerrank 2022-06-15 task (GO) Components in a graph.txt

/*
Components in a graph

https://www.hackerrank.com/challenges/components-in-graph/problem

Report the min (>1) and max sizes of connected subsets.

Sample Input

STDIN   Function
-----   --------
5       bg[] size n = 5
1 6     bg = [[1, 6],[2, 7], [3, 8], [4,9], [2, 6]]
2 7
3 8
4 9
2 6

Sample Output

2 4

This is quite like the 'Merging Communities' I just solved.

Before I read the discussions about "disjoint-set forests",
I just solved it as an array of arrays, a slice of slices.
Do it again, faster and terser:
09:13 AM read problem
09:32 finished stdin, atoi, have to go...
12:12 resume work
12:53 first submit
oops, oops, ouch...
13:33 final submit
*/

package main

import (
	"fmt"
	"io/ioutil"
	// "log"
	"os"
)

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic("ioutil.ReadAll")
	}
	// log.Println(len(bytes))
	arr := make([]int, 0)
	var accu int
	var inDigits bool
	for _, b := range bytes {
		b -= 48
		if b < 10 {
			inDigits = true
			accu = accu*10 + int(b)
		} else {
			if inDigits {
				arr = append(arr, accu)
				accu = 0
			}
			inDigits = false
		}
	}
	if inDigits {
		// at EOF w/o NL
		arr = append(arr, accu)
	}
	// log.Println(arr)

	// make a slice of slices to hold the forest
	n := arr[0]
	ss := make([][]int, 1+2*n) // slot [0] ignored, value range = 1..2n

	// init all members with their own index number, as their group number
	for i := 1; i <= 2*n; i++ {
		ss[i] = []int{i}
	}
	// log.Println(ss)

	// take input pairs to cluster members
	for i := 1; i <= 2*n; i += 2 {
		j := arr[i]
		k := arr[i+1]
		// log.Println("Join", j, k)
		if j == k {
			continue // moot here; for generality
		}
		// What # is the representative member,
		// the leader of each of their groups?
		jrep := ss[j][0]
		krep := ss[k][0]

		if jrep == krep {
			continue // groups already joined
		}
		// could optimize to move shorter list,
		// double code not worth the work here.
		// Just let J inherit all of K members.

		// Append the krep member list to jrep.
		ss[jrep] = append(ss[jrep], ss[krep]...)

		// Update all the krep members to jrep.
		for m := 0; m < len(ss[krep]); m++ {
			ss[ss[krep][m]][0] = jrep
		}

		// Drop member list, if any, from krep.
		ss[krep] = ss[krep][:1]

		// log.Println(ss)
	}
	// find min (>1) and max length slice in ss.
	// math.MaxInt isn't undefined on my computer!
	// but HackerRank says: undefined: math.MaxInt
	min := 1 << 30
	max := 1
	for i := 1; i <= n; i++ {
		len := len(ss[i])
		if min > len && len > 1 {
			min = len
		}
		if max < len {
			max = len
		}
	}
	fmt.Println(min, max)
}
