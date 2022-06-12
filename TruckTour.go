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
 * Complete the 'truckTour' function below.
 *
 * The function is expected to return an INTEGER.
 * The function accepts 2D_INTEGER_ARRAY petrolpumps as parameter.
 */

func truckTour(petrolpumps [][]int32) int32 {
	// Write your code here

	// I see the net benefit per pump is (liters at[1] minus petrol at[0]).
	// A max benefit pump might be a good starting place, try those first(?).
	// After finding any one good starting place,
	// can I change strategy based on one such ring of net residues?
	// Or just starting with the list/ring of benefits, not yet residues,
	// Can I note the most negative as the stumblingblock, and...what?
	// perhaps work backwards summing benefits until overcome that s.b.
	// sample:
	// 3
	// 1 5    benefit = -4 = Step 1: most negative
	// 10 3 benefit = +7 = Step 3: raises -3 to +4, sufficient to pass test (or not?)
	// 3 4    benefit = +1 = Step 2: raises -4 to -3
	//
	// Next problem is that any one benefit is not a net residue.
	// A run of small negatives may be worse than such max neg.
	// I could group successive negs into one big neg, including any around the ring.
	// I could group successive poss into one big pos, as a linear list, not around the ring.
	// I could group alternating {pos,neg} pairs into new 'pump-groups' with signed benefit.
	// Any initial negs cannot satisfy, so lump them with any final negs at end of linear list.
	// There might be no negs in data, but there must be some pos benefit to make a solution
	// Except when N==1, when net benefit could be zero. (In fact even N>0 all benefits==0.)
	// In fact any "benefits" above might be == 0. So my 'pos' s/b non-neg, >= 0.
	// But why lump a just-zero with a next negative, not group with prior neg?
	// Suggests my alg. may be suboptimal...?
	// No, {non-neg, neg} remains sufficient; (Or even just one {non-neg} group.)
	// But it tells me I must repeatedly move any initial neg 'pump-groups' to end.
	// So loop is over when there is only one group.

	// Now, unstead of multiple passes re-grouping groups,
	// Can I generate a one-pass algorithm from start to end?
	// Yes:
	// Any negatives encountered sum into a post-dataArray deficiency accumulator.
	// At start of any Non-Negative sequence, note the starting point as potential answer.
	// If any later negatives, looking one-by-one, wipe out (i.e., <0) the net residual benefit,
	// then discard the potential answer, sum net deficiency to then end, resume looking.

	var ans int32

	// How big? 10^9 x 2 x 10^5 = 10^15 = 1,000,000,000,000,000 = 3 8D7E A4C6 8000 hex.
	// Hmmm. Turns out I don't even need the finalLack, if all test cases are valid.
	// var finalLack int64
	var netResidue int64

	for i := 0; i < len(petrolpumps); i++ {
		netResidue += int64(petrolpumps[i][0] - petrolpumps[i][1])
		if netResidue < 0 {
			// finalLack += netResidue;
			netResidue = 0
			ans = -1
		} else {
			if ans == -1 {
				ans = int32(i)
			}
		}
	}

	return ans
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	nTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	n := int32(nTemp)

	var petrolpumps [][]int32
	for i := 0; i < int(n); i++ {
		petrolpumpsRowTemp := strings.Split(strings.TrimRight(readLine(reader), " \t\r\n"), " ")

		var petrolpumpsRow []int32
		for _, petrolpumpsRowItem := range petrolpumpsRowTemp {
			petrolpumpsItemTemp, err := strconv.ParseInt(petrolpumpsRowItem, 10, 64)
			checkError(err)
			petrolpumpsItem := int32(petrolpumpsItemTemp)
			petrolpumpsRow = append(petrolpumpsRow, petrolpumpsItem)
		}

		if len(petrolpumpsRow) != 2 {
			panic("Bad input")
		}

		petrolpumps = append(petrolpumps, petrolpumpsRow)
	}

	result := truckTour(petrolpumps)

	fmt.Fprintf(writer, "%d\n", result)

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
