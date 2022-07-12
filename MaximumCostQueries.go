// Hackerrank 2022-06-27 task (GO) Super Maximum Cost Queries.txt
// https://www.hackerrank.com/challenges/maximum-cost-queries/problem

// In Go, golang, passes 100% of test cases!
// Considered Hard, only submitted by ~2000.

/*
Where Cost of a path is the MAX edge weight along path,
Report count of paths in [L..R] inclusive weight range.

Another Disjoint Set Union problem.

Idea: progressively grow a DSU from smallest to largest weight,
and sum the various (n * n-1 / 2) paths in all the DSU subsets.

// Suppose I don't re-sum all DS, but only the two
// affected: union len(m) with len(n) => len(m+n).

So, on the sample data:
5 5
1 2 3
1 4 2
2 5 6
3 4 1
1 1
1 2
2 3
2 5
1 6

Resort {u, v, w} by ascending weights w:
3 4 1
1 4 2
1 2 3
2 5 6

after adding w=1: 3-4, there is one DS of 2 nodes, whose edge weight==max weight== 1.
after adding w=2: 1-4, there is one DS of 3 nodes. 3 * 2 / 2 = 3 paths, of max == 2.
after adding w=3: 1-2, there is one DS of 4 nodes. 4 * 3 / 2 = 6 paths, of max == 3.
after adding w=6: 2-6, there is one DS of 5 nodes. 5 * 4 / 2 = 10 paths, of max == 6.

Now, to answer [L-R] inclusive queries:
1-1 => 1 path available at max=1
1-2 => 3 paths available at max=2, and none go away for being below L=1
2-3 => 6 paths available at max=3, minus 1 path available at max=1, = 5
2-5 => ditto, 6 minus 1 = 5
1-6 => all 10 paths available at max=6, and none go away for being below L=1
Which results all agree with desired sample output.

Constraints:
1 <= N,Q <= 10^5
1 <= U,V <= N
1 <= W <= 10^9
1 <= L <= R <= 10^9

1<<30 = 1,073,741,824
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	start := time.Now()

	// stateful data input vars
	var state int
	var n int
	var n3m3 int
	var in int
	var q int
	var q2 int
	var iq int

	// union vars
	var u int
	var v int
	var w int

	// query vars
	var l int
	var r int

	// Disjoint Set Union [1..n] holds nodes.
	// Node n holds a slice of representative
	// AKA parent index, then followers of n.
	var dsu [][]int

	// this to postpone operating on edges until all are input.
	// weights w index map to find a slice of pairs of ints u,v
	var wuv = make(map[int][][2]int)

	// The sum of all DS subset paths grows as edges are added.
	var paths int

	// Ascending list of weight plateaus and their path counts.
	// When added wp[i]=weight, finished acheiving pc[i]=paths:
	var wp = make([]int, 0)
	var pc = make([]int, 0)

	// It turns out to feel very natural to have all the funcs
	// that main calls live inside main to close over its vars.

	// This function performs a DSU union when adding one edge:

	union := func(u, v int) {
		if u == v {
			panic("u==v")
		}
		pu := dsu[u][0]
		pv := dsu[v][0]
		if pu == pv {
			panic("cyclic")
		}
		nu := len(dsu[pu])
		nv := len(dsu[pv])
		if nu < nv {
			// swap
			u, v = v, u
			pu, pv = pv, pu
			nu, nv = nv, nu
		}
		// nu is not less than nv
		// attach pv to bigish pu
		dsu[pu] = append(dsu[pu], dsu[pv]...)
		for _, iv := range dsu[pv] {
			dsu[iv][0] = pu
		}
		dsu[pv] = dsu[pv][:1]
		// fmt.Printf("union %v & %v :: %#v\r\n", u, v, dsu)
		// Now compute the impact on paths
		paths -= (nu - 1) * nu / 2
		paths -= (nv - 1) * nv / 2
		paths += (nu + nv - 1) * (nu + nv) / 2
		// fmt.Printf("union %v & %v :: %v\r\n", u, v, paths)
	}

	// *ALL* my results were slightly low on next test case 2.
	// Some OBO etc... // These two sandwich runAllEdges call:

	sanityCheck1 := func() {
		// verify that wuv contains N-1 edges.
		// sudden Aha! I predict it will not...
		// my map of slices s/b map of references to slices!
		// so any duplicate weight edge would not get saved.
		ne := 0
		for _, suv := range wuv {
			ne += len(suv) // slice of [u,v] pairs at weight w
		}
		// fmt.Println("n", n, "ne", ne) // Ta-Da: n 959 ne 957

		if ne != n-1 {
			panic("ne != n-1")
		}
		// There was just one duplicate in test case 2:
		// Now I caught it: 88888:[[541 528] [916 373]]

	}

	sanityCheck2 := func() {
		// not needed now.
	}

	// This function finishes inputing phase one to build tree
	// after having sorted all edges by ascending edge weight:
	// It builds the ascending wp & pc list to answer queries.

	runAllEdges := func() {

		// process wuv[w] in ascending key w order
		keys := make([]int, 0)
		for w, _ := range wuv {
			keys = append(keys, w)
		}
		// fmt.Println("wuv", wuv)
		sort.Ints(keys)

		// before all real weights,
		// add a negative sentinel
		wp = append(wp, -1)    // -1 is below all non-negative weight plateaus
		pc = append(pc, paths) // path count is initially zero, will grow...

		for _, w := range keys {
			// fmt.Println("run w", w)
			// At weight w, union all its (u,v)[i].
			for i := 0; i < len(wuv[w]); i++ {
				union(wuv[w][i][0], wuv[w][i][1])
			}
			// having finished this w,
			// record path counts for
			// L-R ranges including w.
			wp = append(wp, w)
			pc = append(pc, paths)
			// fmt.Println("w", w, "paths", paths)
		}

		// after all real weights,
		// add a ~maxint sentinel
		wp = append(wp, 1<<30)
		pc = append(pc, paths)

		// fmt.Println("wp", wp)
		// fmt.Println("pc", pc)

		// on the given sample data:
		// wp [-1 1 2 3 6 1073741824]
		// pc: [0 1 3 6 10 10]

	}

	// This function is inputing phase one to build tree.
	// Add input UVW to a map[w]->list(u,v) to run later.

	addEdgeUVW := func() {

		// Note this func inside main closes over u,v,w:
		// _ = u
		// _ = v
		// _ = w

		sp, ok := wuv[w]
		if ok {
			// Finally, here was my downfall:
			// Assigning back to the copy sp,
			// did not add to original wuv[w]
			// WRONG: sp = append(sp, [2]int{u, v})
			wuv[w] = append(sp, [2]int{u, v})
		} else {
			wuv[w] = append(make([][2]int, 0), [2]int{u, v})
		}
		// fmt.Printf("%#v\r\n", wuv)
	}

	// This function is inputing phase two to answer queries.

	answerLR := func() {

		// The main stdin parsing loop set l and r values for me:
		_ = l
		_ = r
		// fmt.Println("\r\nQuerying L,R:", l, r)

		// My original very-mind-troubling plan was to do
		// a binary search in wp for both of l and r, and
		// then interpolate for present versus missing wp.
		// It produced answers that were all slightly low.

		// New plan is to store an entire range of 10^5 w.
		// Today's new approach is invalid: 10^9 not 10^5.
		// Where did that old code go? Ahhh, found it!....

		// N.B. L test uses ">="
		i := sort.Search(len(wp), func(s int) bool { return wp[s] >= l })
		// fmt.Println("searching L=", l, "gave i=", i, "wp[i]=", wp[i])
		// fmt.Println("Subtracting at i-1=", i-1, "wp[i-1]=", wp[i-1], "pc[i-1]=", pc[i-1])

		// N.B. R test uses ">"
		j := sort.Search(len(wp), func(s int) bool { return wp[s] > r })
		// fmt.Println("searching R=", r, "gave j=", j, "wp[j]=", wp[j])
		// fmt.Println("Adding at j-1=", j-1, "wp[j-1]=", wp[j-1], "pc[j-1]=", pc[j-1])

		fmt.Println(pc[j-1] - pc[i-1]) // "ANSWER:"
	}

	// Atoi accumulator
	var accu int

	// This function uses the next input number from stdin parse.
	// It statefully processes n, q, then n triplets, q doublets.

	useAccu := func() {
		// fmt.Println("use", accu)
		switch state {
		case 0:
			// initial state, reading n
			n = accu
			// fmt.Println("n", n)
			n3m3 = 3*n - 3 // N nodes means N-1 edges
			{
				// initialize for n nodes
				dsu = make([][]int, 1+n) // 1-up, slot 0 unused
				for i := 1; i <= n; i++ {
					dsu[i] = []int{i}
				}
			}
			state++
			break
		case 1:
			// second state, reading q
			q = accu
			// fmt.Println("q", q)
			q2 = 2 * q
			state++
			break
		case 2:
			// third state, reading n x triplets
			switch in % 3 {
			case 0:
				u = accu
				// fmt.Println("u", u)
				break
			case 1:
				v = accu
				// fmt.Println("v", v)
				break
			case 2:
				w = accu
				// fmt.Println("w", w)
				addEdgeUVW()
				break
			}
			in++
			if in == n3m3 {
				// before advancing state,
				// process those n inputs.
				sanityCheck1()
				runAllEdges()
				sanityCheck2()
				state++
			}
			break
		case 3:
			// third state, reading q x doublets
			switch iq % 2 {
			case 0:
				l = accu
				// fmt.Println("l", l)
				break
			case 1:
				r = accu
				// fmt.Println("r", r)
				answerLR()
				break
			}
			iq++
			if iq == q2 {
				// Just anal I guess.
				// fmt.Println("EOF")
				state++
			}
			break
		case 4:
		default:
			panic(state)
			break
		}
	}

	// Finally, the main DIY stdin, atoi loop here:

	sb, err := ioutil.ReadAll(os.Stdin)
	check(err)
	var inNum bool
	for _, b := range sb {
		if b < byte('0') {
			// whitespace
			if inNum {
				useAccu()
				accu = 0
				inNum = false
			}
		} else {
			// digit
			accu = accu*10 + int(b-byte('0'))
			inNum = true
		}
	}
	// At EOF w/o NL, a HackerRank trait:
	if inNum {
		useAccu()
	}

	_ = start
	// 7 ms on ~ 1K, 1K data: fmt.Println(time.Since(start).Nanoseconds(), "ns")
}
