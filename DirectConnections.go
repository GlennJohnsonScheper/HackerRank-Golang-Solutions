// DirectConnections.go

// https://www.hackerrank.com/challenges/direct-connections/problem

// ==============================================
// Q: Is this only hard because the main template
// appears wrong: It only reads one line, not two?
// ==============================================
// NP -- I like to DIY.

// Top line t = number of tests.
// then t triplets of lines: n; n * xi coords, n * pi populations
// t in 1..20
// n in 1..100K
// p in 1..10K
// x in 0..10^9 km
// return t answers each mod 1,000,000,007, which t answers are:
// The total km of cable to wire every city to every other city,
// with as many cables as largest of the two cities populations.

/*
sample input:
2
3
1 3 6
10 20 30
5
5 55 555 55555 555555
3333 333 333 33 35

sample output:
280
463055586

plan:
Obviously, sort by descending population.
But that has up to 200K x 200K / 2 loops.
Is there any way to not do so many loops?
Got sample code right. See if performant.
Performant enought. Four fails, negative!
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func main() {
	sb, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	accu := 0
	indigit := false
	for _, b := range sb {
		if b < byte('0') {
			if indigit {
				use(accu)
				accu = 0
				indigit = false
			}
		} else {
			indigit = true
			accu = accu*10 + int(b-byte('0'))
		}
	}
	if indigit {
		use(accu)
	}
}

var state = 0
var t = 0
var ti = 0
var n = 0
var ni = 0

var xp [][]int

func use(val int) {
	switch state {
	case 0:
		t = val
		ti = t
		state++
	case 1:
		n = val
		ni = n
		xp = make([][]int, 0, n)
		state++
	case 2:
		xp = append(xp, []int{val, 0})
		ni--
		if ni == 0 {
			state++
		}
	case 3:
		xp[ni][1] = val
		ni++
		if ni == n {
			solve()
			state++
			ti--
			if ti > 0 {
				state = 1
			}
		}
	default:
		panic("state")
	}
}

func solve() {
	sort.Slice(xp,
		func(i, j int) bool {
			return xp[i][1] > xp[j][1]
		})
	km := int64(0)
	for i := 0; i < len(xp); i++ {
		for j := i + 1; j < len(xp); j++ {
			distance := xp[i][0] - xp[j][0]
			if distance < 0 {
				distance = -distance
			}
			largestP := xp[i][1]
			km += int64(largestP) * int64(distance)
		}
		// periodically knock it down, not at end.
		km = km % 1000000007
	}
	fmt.Println(km)
}
