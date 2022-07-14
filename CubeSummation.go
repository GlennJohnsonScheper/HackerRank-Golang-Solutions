// HackerRank Cube Summation problem

// Sufficient ingenuity if just loop over the few < 1000 updates
// rather then over the < 100^3 cells of the cube for each query.

// But, the provided template retn int32 had to be fixed to int64.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
 * Complete the 'cubeSum' function below.
 *
 * The function is expected to return an INTEGER_ARRAY.
 * The function accepts following parameters:
 *  1. INTEGER n
 *  2. STRING_ARRAY operations
 */

func cubeSum(n int32, operations []string) []int64 {
	// Write your code here
	ans := []int64{}
	// Hold the sparse cube as a slice of update quadruples.
	// Oops. Not quite. New x,y,z must update, not append w.
	// So change slice to map of x,y,z triples looking up w.
	utw := make(map[[3]int]int)
	for _, s := range operations {
		ss := strings.Split(s, " ")
		switch len(ss) {
		case 5:
			// update
			ut := [3]int{}
			ut[0], _ = strconv.Atoi(ss[1])
			ut[1], _ = strconv.Atoi(ss[2])
			ut[2], _ = strconv.Atoi(ss[3])
			w, _ := strconv.Atoi(ss[4])
			utw[ut] = w
			break
		case 7:
			// query
			x1, _ := strconv.Atoi(ss[1])
			y1, _ := strconv.Atoi(ss[2])
			z1, _ := strconv.Atoi(ss[3])
			x2, _ := strconv.Atoi(ss[4])
			y2, _ := strconv.Atoi(ss[5])
			z2, _ := strconv.Atoi(ss[6])
			sum := int64(0)
			for ut, w := range utw {
				if x1 <= ut[0] && x2 >= ut[0] &&
					y1 <= ut[1] && y2 >= ut[1] &&
					z1 <= ut[2] && z2 >= ut[2] {
					sum += int64(w)
				}
			}
			ans = append(ans, sum)
			break
		default:
			panic("not 5,7")
			break
		}
	}
	// fmt.Println("ans", ans)
	return ans
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	TTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	T := int32(TTemp)

	for TItr := 0; TItr < int(T); TItr++ {
		firstMultipleInput := strings.Split(strings.TrimSpace(readLine(reader)), " ")

		matSizeTemp, err := strconv.ParseInt(firstMultipleInput[0], 10, 64)
		checkError(err)
		matSize := int32(matSizeTemp)

		mTemp, err := strconv.ParseInt(firstMultipleInput[1], 10, 64)
		checkError(err)
		m := int32(mTemp)

		var ops []string

		for i := 0; i < int(m); i++ {
			opsItem := readLine(reader)
			ops = append(ops, opsItem)
		}

		res := cubeSum(matSize, ops)

		// fmt.Println("res", res)

		for i, resItem := range res {
			fmt.Fprintf(writer, "%d", resItem)

			if i != len(res)-1 {
				fmt.Fprintf(writer, "\n")
			}
		}

		fmt.Fprintf(writer, "\n")
	}

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
