// Create a Contacts application with the two basic operations: add and find.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	// "time"
)

// Trie as dynamic array, grows by 28 ints.
// First [0-25] are indices of child nodes.
// Final [26] is a count, serves as a bool.
// Naive Trie was slow. Let [27] cache sum.

type Trie []int

var int28 = make([]int, 28)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// start := time.Now()
	// DIY stdio, parse...
	ba, err := ioutil.ReadAll(os.Stdin)
	check(err)

	reCRLFs := regexp.MustCompile(`[[:space:]]+`)
	sa := reCRLFs.Split(string(ba), -1)
	// fmt.Println(len(sa))

	n, err := strconv.Atoi(sa[0])
	check(err)
	if n*2+1 != len(sa) {
		panic("input")
	}

	// fmt.Println("input finished", time.Since(start).Milliseconds())
	// start = time.Now()

	t := Trie{}
	t = append(t, int28...)

	for i := 2; i < len(sa); i += 2 {
		switch sa[i-1] {
		case "add":
			t.add(sa[i])
			break

		case "find":
			t.find(sa[i])
			break

		default:
			panic("cmd")
			break
		}
	}
	// fmt.Println("process finished", time.Since(start).Milliseconds())
}

func (t *Trie) add(s string) {
	// fmt.Println("add", s)
	i := 0
	for _, c := range s {
		(*t)[i+27]++ // count this new child being added somewhere below
		o := int(c) - int('a')
		j := (*t)[i+o]
		if j == 0 {
			j = len(*t)               // offset to new node
			*t = append(*t, int28...) // new node
			(*t)[i+o] = j
		}
		i = j
	}
	(*t)[i+26]++ // count-bool, found
	(*t)[i+27]++ // sum of self,now + any children later

	/*
		for i = 0; i < len(*t); i += 28 {
			// fmt.Println((*t)[i : i+28])
		}
	*/
}

func (t *Trie) find(s string) {
	// fmt.Println("find", s)
	i := 0
	for _, c := range s {
		o := int(c) - int('a')
		j := (*t)[i+o]
		if j == 0 {
			// entire prefix was not found
			fmt.Println("0")
			return
		}
		i = j
	}
	// prefix was found

	// fast method:
	fmt.Println((*t)[i+27])

	/****************************
	 * this was the naive method:
	 * way too slow!
	 *
	// need to sum up self and children
	sum := 0
	queue := []int{}
	queue = append(queue, i)
	for len(queue) > 0 {
		j := queue[0]
		queue = queue[1:]
		sum += (*t)[j+26]
		for k := 0; k < 26; k++ {
			if (*t)[j+k] > 0 {
				queue = append(queue, (*t)[j+k])
			}
		}
	}
	fmt.Println(sum)
	 ****************************/
}
