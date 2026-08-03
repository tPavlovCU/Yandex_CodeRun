package main

import (
	"bufio"
	"os"
	"sort"
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

	return res * negative
}

func writeInt(writer *bufio.Writer, n int) {
	if n == 0 {
		writer.WriteByte('0')
		return
	}

	var buf [20]byte
	pos := len(buf)

	for n > 0 {
		pos--
		buf[pos] = byte(n%10) + '0'
		n /= 10
	}

	writer.Write(buf[pos:])
}

type tower struct {
	start int
	end   int
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	cities := readInt(reader)
	towers := readInt(reader)

	towersSlice := make([]tower, towers)

	for i := 0; i < towers; i++ {
		l := readInt(reader)
		r := readInt(reader)
		towersSlice[i] = tower{l, r}
	}

	sort.Slice(towersSlice, func(i, j int) bool {
		if towersSlice[i].start == towersSlice[j].start {
			return towersSlice[i].end > towersSlice[j].end
		}
		return towersSlice[i].start < towersSlice[j].start
	})

	cnt := 0
	idx := 0
	pos := 0

	for pos < cities {
		furthest := pos

		for idx < towers && towersSlice[idx].start <= pos+1 {
			if towersSlice[idx].end > furthest {
				furthest = towersSlice[idx].end
			}
			idx++
		}

		if furthest == pos {
			writer.WriteString("No")
			return
		}

		cnt++
		pos = furthest
	}

	writer.WriteString("Yes\n")
	writeInt(writer, cnt)
	writer.WriteByte('\n')
}
