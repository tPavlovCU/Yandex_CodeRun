package main

import (
	"bufio"
	"os"
	"strings"
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
func writeInt(writer *bufio.Writer, n int64) {
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
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()
	_ = reader

	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	var result [26]int

	cnt := int64(0)
	lineB := []byte(line)
	for _, value := range lineB {
		idx := value - 'a'
		if idx < 25 {
			for m := idx + 1; m < 26; m++ {
				cnt += int64(result[m])
			}
		}
		result[idx] += 1
	}
	writeInt(writer, cnt)
}
