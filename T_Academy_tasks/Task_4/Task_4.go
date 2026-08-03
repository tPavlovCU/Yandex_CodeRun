package main

import (
	"bufio"
	"os"
)

func readInt(reader *bufio.Reader) int {
	var res int
	b, err := reader.ReadByte()
	negative := 1
	for err == nil && (b == '\n' || b == ' ' || b == '\r') {
		b, err = reader.ReadByte()
	}

	if err != nil {
		return 0
	}

	if b == '-' {
		negative = -1
		b, err = reader.ReadByte()
	}

	for err == nil && b >= '0' && b <= '9' {
		res = res*10 + int(b-'0')
		b, err = reader.ReadByte()
	}
	res = res * negative
	return res
}
func writeInt(writer *bufio.Writer, n int) {
	if n == 0 {
		writer.WriteByte('0')
		return
	}
	negative := false

	var buf [20]byte
	if n < 0 {
		negative = true
		n = -1 * n
	}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte((n % 10) + '0')
		n /= 10
	}
	if negative {
		writer.WriteByte('-')
	}
	writer.Write(buf[pos:])
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minLenSum(line []int, badSum int) int {
	mapa := make(map[int]int, len(line))
	result := len(line) + 1
	sum := 0
	mapa[0] = -1
	for idx, value := range line {
		sum += value
		valueOk, ok := mapa[sum-badSum]
		if ok {
			result = min(result, idx-valueOk)
		}
		mapa[sum] = idx
	}
	if sum == 0 {
		return 0
	}
	if result == len(line)+1 {
		return -1
	}
	return result
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()
	_ = reader
	N := readInt(reader)
	commands := make([]int, N)
	lCnt := 0
	rCnt := 0
	for n := 0; n < N; n++ {
		symbol, _ := reader.ReadByte()
		if symbol == 'R' {
			commands[n] = 1
			rCnt += 1
		} else if symbol == 'L' {
			commands[n] = -1
			lCnt += 1
		}
	}
	res := minLenSum(commands, rCnt-lCnt)
	writeInt(writer, res)
}
