// Hackerrank 2022-06-12 task (GO) Merging Communities.txt

/*
11:30 am read, started thinking about it.
12:30 pm got the strategy, ready to code.
2:30 pm passes sample data. deinstrument.
6/10 fails. Spend 5 hackos. reinstrument.
TMI in dumps! Add a sanity check to code.
3:30 Stopped making an as..., two places.

read N and Q.
n in 1...10^5
q in 1...2*10^5

Immediately, people[1...n] exist each in a group of size 1.
Queries start with 'M' = merge, or 'Q' = query cluster SIZE.

Sample input
3 6 -- n = 3, q = 6
Q 1 -- print the size of the community containing person 1
M 1 2 -- merge the communities containing persons 1 and 2
Q 2
M 2 3
Q 3
Q 2

Sample output
1
2
3
3

I recall some clustering idea of representing all members
by any one representative member of the group's number...

So keep an array of arr[i]-->j, their leader.
A leader must keep the list of group members.
So in two steps, any member can find members.
Joining two groups, rewrite number to leader.

First thought is make an array of slices,
either naming a leader or listing members.
If such arr[i]'th slice[0] is some j != i,
j is his leader, else i is himself leader.
If len(slice) > 1, group has more members.

*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	// fmt.Println("Hello World")

	/* getting my Go legs...
	// this is one empty slice:
	one := make([]int, 0)
	one = append(one, 1)
	one = append(one, 2)
	// this is big slice of slices:
	many := make([][]int, 3)
	many[0] = one
	// fmt.Println(many)
	... That all works */

	// I should stop reading all stdin to a buffer,
	// and learn to use the GO line-by-line idioms.
	// Irks me though, tending towards inefficient.

	sc := bufio.NewScanner(os.Stdin)
	needTopLine := true
	var n int
	var q int
	var person [][]int
	// idiomatic go: loop ends on eof or err:
	for sc.Scan() {
		line := strings.TrimRight(strings.TrimRight(sc.Text(), "\n"), "\r")
		if needTopLine {
			twain := strings.Split(line, " ")
			if len(twain) != 2 {
				panic("top line")
			}
			n, _ = strconv.Atoi(twain[0])
			q, _ = strconv.Atoi(twain[1])
			// fmt.Println(n, q)
			// Generate the slice of N slices.
			// HackerRank's people index 1-up.
			person = make([][]int, n+1)
			for i := 1; i <= n; i++ {
				person[i] = append(make([]int, 0), i)
			}
			// fmt.Println(person)
			needTopLine = false
			continue
		}
		// else not top line
		tokens := strings.Split(line, " ")
		if len(tokens) < 1 {
			panic("empty line")
		}
		switch tokens[0] {
		case "M":
			i, _ := strconv.Atoi(tokens[1])
			j, _ := strconv.Atoi(tokens[2])
			// fmt.Printf("Merge %v, %v\n", i, j)

			// check assumption
			if i == j {
				// THIS OCCURED: panic("i==j")
				// Expletives deleted!
				continue
			}

			iLeader := person[i][0]
			jLeader := person[j][0]

			// check for yet another stumble.
			// Adding test stopped my panics.
			if iLeader == jLeader {
				continue
			}

			// for sanity check
			leniL := len(person[iLeader])
			lenjL := len(person[jLeader])

			if len(person[iLeader]) < len(person[jLeader]) {
				// merge shorter iLeader into longer jLeader
				// append list of followers
				person[jLeader] = append(person[jLeader], person[iLeader]...)
				// tell everyone in iLeader, they now follow jLeader
				for _, k := range person[iLeader] {
					person[k][0] = jLeader
				}
				// in case iLeader had followers, drop them
				person[iLeader] = person[iLeader][:1]

				// sanity check
				if person[i][0] != person[j][0] {
					panic("1st i-j test")
				}
				if person[iLeader][0] != person[jLeader][0] {
					panic("1st iLeader-jLeader test")
				}
				if len(person[i]) != 1 {
					panic("len i not 1 test")
				}
				if len(person[iLeader]) != 1 {
					panic("len iLeader not 1 test")
				}
				if len(person[jLeader]) != leniL+lenjL {
					panic("len jLeader not sum test")
				}

			} else {
				// merge poss. shorter jLeader into poss. longer iLeader
				// append list of followers
				person[iLeader] = append(person[iLeader], person[jLeader]...)
				// tell everyone in jLeader, they now follow iLeader
				for _, k := range person[jLeader] {
					person[k][0] = iLeader
				}
				// in case jLeader had followers, drop them
				person[jLeader] = person[jLeader][:1]

				// sanity check
				if person[i][0] != person[j][0] {
					panic("2nd i-j test")
				}
				if person[iLeader][0] != person[jLeader][0] {
					panic("2nd iLeader-jLeader test")
				}
				if len(person[j]) != 1 {
					panic("len j not 1 test")
				}
				if len(person[jLeader]) != 1 {
					panic("len jLeader not 1 test")
				}
				if len(person[iLeader]) != leniL+lenjL {
					panic("len iLeader not sum test")
				}

			}
			break
		case "Q":
			i, _ := strconv.Atoi(tokens[1])
			// fmt.Printf("Report on %v\n", i)
			iLeader := person[i][0]
			fmt.Println(len(person[iLeader])) // desired answer
			break
		default:
			panic("not M, not Q")
			break
		}
		// way way TMI -- // fmt.Println(person)
		_ = q
	}
	if err := sc.Err(); err != nil && err != io.EOF {
		panic(err)
	}
	// fmt.Println("fini")
}
