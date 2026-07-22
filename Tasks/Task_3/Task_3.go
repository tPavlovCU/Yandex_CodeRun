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

	moves := make([][]byte, N)
	for i := int32(0); i < N; i++ {
		moves[i] = make([]byte, M)
	}

	for m := int32(0); m < M; m++ {
		item := readInt(reader)
		if m == 0 {
			ar[m] = item
		} else {
			ar[m] = ar[m-1] + item
			moves[0][m] = 'R'
		}
	}

	for n := int32(1); n < N; n++ {

		item0 := readInt(reader)
		ar[0] = ar[0] + item0
		moves[n][0] = 'D'

		for m := int32(1); m < M; m++ {
			item := readInt(reader)

			if ar[m-1] > ar[m] {
				ar[m] = ar[m-1] + item
				moves[n][m] = 'R'
			} else {
				ar[m] = ar[m] + item
				moves[n][m] = 'D'
			}
		}
	}

	writeInt32(writer, ar[M-1])
	writer.WriteByte('\n')

	pathBuf := make([]byte, (N-1)+(M-1))
	pathPos := len(pathBuf)

	currN := N - 1
	currM := M - 1

	for currN > 0 || currM > 0 {
		pathPos--
		step := moves[currN][currM]
		pathBuf[pathPos] = step

		if step == 'R' {
			currM--
		} else {
			currN--
		}
	}

	for i := pathPos; i < len(pathBuf); i++ {
		writer.WriteByte(pathBuf[i])
		if i < len(pathBuf)-1 {
			writer.WriteByte(' ') // Добавляем пробел между шагами
		}
	}
	writer.WriteByte('\n')
}
