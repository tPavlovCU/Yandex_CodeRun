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
	result := make([][]int32, 0, N)
	for n := int32(0); n < N; n++ {
		new := make([]int32, M)
		result = append(result, new)
	}

	for n := int32(0); n < N; n++ {
		result[n][0] = 0
		for m := int32(0); m < M; m++ {
			flag1 := (n < 2 || m < 1)
			flag2 := (n < 1 || m < 2)
			if flag1 && flag2 {
				result[n][m] = 0
			} else if flag1 {
				result[n][m] = result[n-1][m-2]
			} else if flag2 {
				result[n][m] = result[n-2][m-1]
			} else {
				result[n][m] = result[n-2][m-1] + result[n-1][m-2]
			}
		}
		result[0][0] = 1
	}

	writeInt32(writer, result[N-1][M-1])
}
