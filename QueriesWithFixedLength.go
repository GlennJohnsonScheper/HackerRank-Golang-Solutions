package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
 * Complete the 'solve' function below.
 *
 * The function is expected to return an INTEGER_ARRAY.
 * The function accepts following parameters:
 *  1. INTEGER_ARRAY arr
 *  2. INTEGER_ARRAY queries
 */

func solve(arr []int32, queries []int32) []int32 {
	// Write your code here

	// This all new algorithm cut the time to compute
	// HR testcase 3 down from 500+ seconds to 78 ms.

	// log.Printf("Array = [%v]\n", arr)
	// log.Printf("queries = [%v]\n", queries)

	// Here store least-max answer for each query value.
	ans := make([]int32, len(queries))

	// My first algorithm took 500 seconds on 100K array.
	// Now I think to amortize sorting among the queries.
	// Is there even time to sort 100K array? Int: 24 ms.
	// Now put {value, index} slice to sort. Still 24 ms!
	vis := make([]Pairs, len(arr))
	for i, v := range arr {
		vis[i].value = v
		vis[i].index = int32(i)
	}
	sort.Slice(vis, func(i, j int) bool {
		// when anonymous "less" func uses '<', it would sort ascending.
		// because I reversed it here to '>', this will sort descending.
		return vis[i].value > vis[j].value
	})

	// With a sorted vis beside arr to use on each query,
	// I think tallest values must block out adjacencies
	// from having a less tall max if they include index.
	// So use vis indexes in order to split arr at index.
	// Stop when no partition is as long as query length.
	// Perhaps do a binary tree of {lsize, index, rsize}.

	// Run each query length value
	for qnum, qlen := range queries {
		// log.Println("Running qlen:", qlen, "at qnum:", qnum)

		// Inject these loop-control accounting details into tree
		accounting := Accounting{
			inputArrayLength: int32(len(arr)),
			qlen:             qlen,
			unspoiled:        int32(len(arr)),
		}

		// At this qlen[qnum], start a new tree of nodes.
		var arrPartitions Tree

		// Start using the tallest array values first
		for _, pair := range vis {
			// log.Printf("Do tallest value %v at index %v\n", pair.value, pair.index)
			arrPartitions.insert(pair.index, &accounting)
			if accounting.unspoiled < qlen {
				// log.Printf("Stopping as unspoiled (%v) < qlen (%v)\n", accounting.unspoiled, qlen)
				ans[qnum] = pair.value
				break
			}
		}
		// show me!
		// printInOrder(arrPartitions.root, 0)
		// add more diagnostics as the first loop's answer came out wrong:
		// I got it now -- break
	}
	// log.Printf("answer = [%v]\n", ans)
	return ans
}

// This struct is to sort the values, but remember their indexes
type Pairs struct {
	value int32
	index int32
}

// This binary search tree was straight from a
// googled bogotobogo.com BST golang tutorial:
type Tree struct {
	root *Node
}

// I had delegated loop math and logic into the tree,
// but separate that concern from it now for clarity.
type Accounting struct {
	// Tree needs inputArrayLength to solve top lside & rside.
	inputArrayLength int32
	// Tree needs qlen to know if a sub-partition is unusable.
	qlen int32
	// Tree needs unspoiled to deduct unusable partition size.
	// But the caller himself can test unspoiled to stop loop.
	unspoiled int32
}

type Node struct {
	left  *Node
	lsize int32 // how many ints are left of this index
	key   int32 // these tree keys being my indexes in arr
	rsize int32 // how many ints are right of this index
	right *Node
}

// Tree .insert() = non-recursive first/top call only
func (t *Tree) insert(index int32, accounting *Accounting) {
	// passed index ranges from [0 to accounting.inputArrayLength - 1]
	if t.root == nil {
		// This call will add the root, and there will be no recursion.
		// root node computes lside, rside from inputArrayLength & index.
		// at top node, lsize = index, holding [0 to index-1]
		// at top node, rsize = inputArrayLength - 1 - index, holding [index+1 to len-1]

		// dblchk OBO errors: when min index == 0, ls s/b 0:
		ls := index

		// dblchk OBO errors: when max index == len-1, rs s/b 0:
		rs := accounting.inputArrayLength - 1 - index

		accounting.unspoiled-- // for this one index doing partitioning
		if ls < accounting.qlen {
			// the left side is spoiled, cannot hold any qlen partitions.
			accounting.unspoiled -= ls
		}
		if rs < accounting.qlen {
			// the right side is spoiled, cannot hold any t.qlen partitions.
			accounting.unspoiled -= rs
		}
		t.root = &Node{lsize: ls, key: index, rsize: rs}
	} else {
		// these are all lower nodes, get lside, rside from their parents.
		t.root.insert(index, accounting)
	}
}

// Node .insert() = recursive helper
func (n *Node) insert(index int32, accounting *Accounting) {
	if index <= n.key {
		if n.left == nil {
			// if there is not yet any n.left child,
			// we may repartition the n.lsize bytes,
			// for that is the run containing index.
			// But, n.lsize < qlen: Forget about it!
			// Else I erroneously do an unspoiled--.
			if n.lsize >= accounting.qlen {
				// for a new left node (index is less than n.key),
				// the new node's (closer) rside is n.key - index - 1;
				// the new node's (distal) lside is n.lsize - (n.key - index)

				// dblchk OBO errors: when min index == key-lsize, ls s/b 0:
				// the min index at bottom of the lside is n.key - n.lside.
				// so new space left of index = index - (n.key - n.lside).
				ls := n.lsize - (n.key - index)

				// dblchk OBO errors: when max index == key-1, rs s/b 0:
				rs := n.key - 1 - index

				accounting.unspoiled-- // for this one index doing partitioning
				if ls < accounting.qlen {
					// the left side is spoiled, cannot hold any qlen partitions.
					accounting.unspoiled -= ls
				}
				if rs < accounting.qlen {
					// the right side is spoiled, cannot hold any t.qlen partitions.
					accounting.unspoiled -= rs
				}

				n.left = &Node{lsize: ls, key: index, rsize: rs}
			}

		} else {
			n.left.insert(index, accounting) // left-recursion
		}
	} else {
		if n.right == nil {
			// if there is not yet any n.right child,
			// we may repartition the n.rsize bytes,
			// for that is the run containing index.
			// But, n.rsize < qlen: Forget about it!
			// Else I erroneously do an unspoiled--.
			if n.rsize >= accounting.qlen {
				// for a new right node (index is more than n.key),
				// the new node's (closer) lsize is index - n.key - 1;
				// the new node's (distal) rsize is n.rsize - (index - n.key)

				// dblchk OBO errors: when min index == key+1, ls s/b 0:
				ls := index - n.key - 1

				// dblchk OBO errors: when max index == key+rsize, rs s/b 0:
				// Terse now, as I have stared down a double negation above.
				rs := n.rsize - (index - n.key)

				accounting.unspoiled-- // for this one index doing partitioning
				if ls < accounting.qlen {
					// the left side is spoiled, cannot hold any qlen partitions.
					accounting.unspoiled -= ls
				}
				if rs < accounting.qlen {
					// the right side is spoiled, cannot hold any t.qlen partitions.
					accounting.unspoiled -= rs
				}
				n.right = &Node{lsize: ls, key: index, rsize: rs}
			}
		} else {
			n.right.insert(index, accounting) // right-recursion
		}
	}
}

func printInOrder(n *Node, depth int) {
	if n == nil {
		return
	} else {
		printInOrder(n.left, depth+1)
		// log.Printf("Depth: %v -- lsize: %v [key: %v] rsize: %v\n", depth, n.lsize, n.key, n.rsize)
		printInOrder(n.right, depth+1)
	}
}

// This was the 500+ second correct but too slow version 1:

func too_slow_solve(arr []int32, queries []int32) []int32 {
	// Write your code here
	// fmt.Printf("Array = [%v]\n", arr)
	// fmt.Printf("queries = [%v]\n", queries)

	var ans []int32
	for _, qlen := range queries {
		// run a sortedWindow of len=qlen, across input data of len(arr).
		sortedWindow := make([]int32, qlen)
		// sorted-insert the first qlen data into sortedWindow.
		for i := int32(0); i < qlen; i++ {
			// during filling loop, len(filled) equals i.
			tgt1 := findIndexOf(arr[i], sortedWindow, i)
			// fmt.Printf("for arr[%v]=%v found tgt1 at %v\n", i, arr[i], tgt1)
			// if tgt1 < len(filled), must slide some up.
			for j := i; j > tgt1; j-- {
				sortedWindow[j] = sortedWindow[j-1]
			}
			sortedWindow[tgt1] = arr[i]
			// fmt.Printf("sortedWindow = [%v]\n", sortedWindow)
		}

		// Every maximum is/will be at final sortedWindow[qlen-1].
		// Noting here the least result obtained after each slide:
		least := sortedWindow[qlen-1]

		// sorted-insert the rest of len(arr)-qlen data into sortedWindow.
		// In the same moment, deleting eldest element leaving the window.
		for i := qlen; i < int32(len(arr)); i++ {
			// during sliding loop, len(filled) equals qlen.
			// arr[i] is entering window
			// arr[i-qlen] is exiting window
			tgt1 := findIndexOf(arr[i], sortedWindow, qlen)
			// Spent 5 hackos to see my data out of order!
			// Need a distinction here:
			// sortedWindow[tgt1] MIGHT be a match to value.
			// But more likely, there was NO MATCH FOUND,
			// in which case tgt1 indexes a next larger value.
			// log.Printf("for arr[%v]=%v found tgt1 at %v\n", i, arr[i], tgt1)
			tgt2 := findIndexOf(arr[i-qlen], sortedWindow, qlen)
			// There MUST BE a match to tgt2, about to be removed. Verify:
			if tgt2 == qlen {
				panic("tgt2 was not found")
			}
			// log.Printf("but arr[%v]=%v found tgt2 at %v\n", i-qlen, arr[i-qlen], tgt2)
			// sortedWindow[tgt1] is entering window
			// sortedWindow[tgt2] is exiting window

			// fighting OBO errors...
			// separate these two cases:
			if tgt1 < tgt2 {
				// if tgt1 < tgt2, must slide some up.
				for j := tgt2; j > tgt1; j-- {
					sortedWindow[j] = sortedWindow[j-1]
				}
			}
			if tgt1 > tgt2 {
				// if tgt1 > tgt2, must slide some down.

				// but tgt1 not found reports tgt1=qlen.
				// This fixes a panic on easy test data:
				// if(tgt1 == qlen) {
				//     tgt1 = qlen - 1
				// }
				// No, that was incorrect fix.
				// tgt1 was ALWAYS 1 too high!
				tgt1--
				for j := tgt2; j < tgt1; j++ {
					sortedWindow[j] = sortedWindow[j+1]
				}
			}
			sortedWindow[tgt1] = arr[i]
			// fmt.Printf("sortedWindow = [%v]\n", sortedWindow)

			// Diagnostic: Verify alg. has not mis-sorted the window.
			{
				for j := int32(1); j < qlen; j++ {
					if sortedWindow[j] < sortedWindow[j-1] {
						panic("Disordered sort in Window")
					}
				}
			}

			// Within the newly sorted window, update least
			if least > sortedWindow[qlen-1] {
				least = sortedWindow[qlen-1]
			}
		}
		ans = append(ans, least)
	}
	return ans
}

func findIndexOf(val int32, A []int32, n int32) int32 {
	// straight from the wikipedia Binary_search_algorithm
	L := int32(0)
	R := int32(n - 1)
	for L <= R {
		m := (L + R) / 2
		if A[m] < val {
			L = m + 1
		} else if A[m] > val {
			R = m - 1
		} else {
			return m
		}
	}
	// for unsuccessful outcome, i.e., val not-found,
	// the current L might be in valid [0:n-1] range
	// but for big val, L = past valid [0:n-1] range,
	// perfect for appending big val during filling.
	return L // for unsuccessful
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	firstMultipleInput := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	nTemp, err := strconv.ParseInt(firstMultipleInput[0], 10, 64)
	checkError(err)
	n := int32(nTemp)

	qTemp, err := strconv.ParseInt(firstMultipleInput[1], 10, 64)
	checkError(err)
	q := int32(qTemp)

	arrTemp := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	var arr []int32

	for i := 0; i < int(n); i++ {
		arrItemTemp, err := strconv.ParseInt(arrTemp[i], 10, 64)
		checkError(err)
		arrItem := int32(arrItemTemp)
		arr = append(arr, arrItem)
	}

	var queries []int32

	for i := 0; i < int(q); i++ {
		queriesItemTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
		checkError(err)
		queriesItem := int32(queriesItemTemp)
		queries = append(queries, queriesItem)
	}

	result := solve(arr, queries)

	for i, resultItem := range result {
		fmt.Fprintf(writer, "%d", resultItem)

		if i != len(result)-1 {
			fmt.Fprintf(writer, "\n")
		}
	}

	fmt.Fprintf(writer, "\n")

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
