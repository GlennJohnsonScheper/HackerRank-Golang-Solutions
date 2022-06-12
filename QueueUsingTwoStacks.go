// My first HackerRank challenge in Go, the new Lingua Franca.
// After 2 unlocks, I still have 3 test cases failing timeout.
// Reading discussions, I will now apply @Tharsanan's hot tip:
// Cache a last peek value to not have to {repile twice} * N.

package main

import (
"os"
"bufio"
"strconv"
"strings"
"time"
)

// There is one inputting stack and one outputting stack.
// All the tokens must be left in the one stack last used.
// Before queuing, move all tokens to the inputting stack.
// Before dequeuing, move all tokens to the outputting stack.

// so make a Stack type

type Stack struct {
    // containing a Slice of ints as a dynamic array to hold the Stack
    Slice []int
}

// containing member funcs Queue, Dequeue, Peek, Empty

func (s *Stack) Queue(a int) {
    s.Slice = append(s.Slice, a)
}

func (s *Stack) Dequeue() {
    s.Slice = s.Slice[:len(s.Slice)-1]
}

func (s *Stack) Peek() int {
    return s.Slice[len(s.Slice)-1]
}

func (s *Stack) Empty() bool {
    return len(s.Slice) == 0
}

// GiveItAllUpTo is my first attempt to speed up: rid the Push(Pop()) loop:
// too slow:
// for ; !output.Empty() ; {
//     input.Queue(output.Peek())
//     output.Dequeue()
// }

func (s *Stack) GiveItAllUpTo(t *Stack) {
    if len(s.Slice) == 0 {
        return
    }
    // reverse my data in place
    for i, j := 0, len(s.Slice)-1; i < j; i, j = i+1, j-1 {
        s.Slice[i], s.Slice[j] = s.Slice[j], s.Slice[i]
    }
    // give it away
    t.Slice = append(t.Slice, s.Slice...)
    // empty mine
    s.Slice = nil
}

func main() {
    started := time.Now()
    // First Gigantic lesson in GO is I must do my own buffered IO.
    bufReader := bufio.NewReaderSize(os.Stdin, 16 * 1024 * 1024)
    bufStdout := bufio.NewWriterSize(os.Stdout, 16 * 1024 * 1024)

    // step 1: process just one top line.
    rawLine, _, _ := bufReader.ReadLine()
    lines, _ := strconv.Atoi(strings.TrimRight(strings.TrimRight(string(rawLine), "\n"), "\r"))

    var input Stack
    var output Stack

    // adjust my IO to match revealed test case output, lacking final 0x0A at EOF.
    needLf := false

    lastPeeked := 0
    
    for i := 0; i < lines; i++ {
        rawLine, _, _ := bufReader.ReadLine()
        c := rawLine[0]
        switch rune(c) {
            case '1':
                value, _ := strconv.Atoi(strings.TrimRight(strings.TrimRight(string(rawLine[2:]), "\n"), "\r"))
                output.GiveItAllUpTo(&input)
                if(input.Empty()) {
                    lastPeeked = value;
                }
                input.Queue(value)
                break
            case '2':
                input.GiveItAllUpTo(&output)
                output.Dequeue();
                // Now is the cheap time to refresh lastPeeked
                // Oops, only if non-empty! Else Panic.
                if(!output.Empty()) {
                    lastPeeked = output.Peek();
                }
                break
            case '3':
                // This costly repiling step was saved, *2, *N:
                // On my laptop, on unlocked test case 15 data,
                // Without the next line, one run took 1212 ms.
                // Re-adding .Give* line, one run took 5750 ms.
                // Great thanks to @Tharsanan for that hot tip.
                // input.GiveItAllUpTo(&output)
                // var peeking int
                // peeking = output.Peek();
                if(needLf) {
                    bufStdout.WriteByte(0x0A)
                }
                needLf = true
                bufStdout.WriteString(strconv.Itoa(lastPeeked))
                break
        }
    }
    bufStdout.Flush()
    os.Stdout.Close()
    _ = started
    // os.Stderr.WriteString(strconv.Itoa(int(time.Since(started).Milliseconds())))
}
