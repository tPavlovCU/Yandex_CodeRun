package main

import (
	"bufio"
	"os"
)

func readInt(reader *bufio.Reader) int32 {
	var res int32
	b, err := reader.ReadByte()
	negative := int32(1)
	for err == nil && (b == '\n' || b == ' ' || b == '\r') {
		b, err = reader.ReadByte()
	}

	if err != nil {
		return 0
	}

	if b == '-' {
		negative = int32(-1)
		b, err = reader.ReadByte()
	}

	for err == nil && b >= '0' && b <= '9' {
		res = res*10 + int32(b-'0')
		b, err = reader.ReadByte()
	}
	res = res * negative
	return res
}
func writeInt32(writer *bufio.Writer, n int32) {
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
func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<11)
	writer := bufio.NewWriterSize(os.Stdout, 1<<4)
	defer writer.Flush()

	N := readInt(reader)
	M := readInt(reader)

	ar := make([]int32, M)
	for m := int32(0); m < M; m++ {
		item := readInt(reader)
		if m == 0 {
			ar[m] = item
		} else {
			ar[m] = ar[m-1] + item
		}
	}
	for n := int32(1); n < N; n++ {
		item := readInt(reader)
		ar[0] = ar[0] + item
		for m := int32(1); m < M; m++ {
			item := readInt(reader)
			ar[m] = min(ar[m-1], ar[m]) + item
		}
	}

	writeInt32(writer, ar[M-1])

}
