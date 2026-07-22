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
	line1 := make([]int32, 0, N)
	for n := int32(0); n < N; n++ {
		item := readInt(reader)
		line1 = append(line1, item)
	}

	M := readInt(reader)
	line2 := make([]int32, 0, M)
	for m := int32(0); m < M; m++ {
		item := readInt(reader)
		line2 = append(line2, item)
	}

	ar := make([][]int32, N+1)
	for i := int32(0); i <= N; i++ {
		ar[i] = make([]int32, M+1)
	}

	for p := int32(0); p < N; p++ {
		for k := int32(0); k < M; k++ {
			if line1[p] == line2[k] {
				ar[p+1][k+1] = ar[p][k] + 1
			} else {
				if ar[p][k+1] > ar[p+1][k] {
					ar[p+1][k+1] = ar[p][k+1]
				} else {
					ar[p+1][k+1] = ar[p+1][k]
				}
			}
		}
	}

	ansLength := ar[N][M]
	if ansLength == 0 {
		return
	}

	ans := make([]int32, 0, ansLength)
	i, j := N, M
	for i > 0 && j > 0 {
		if line1[i-1] == line2[j-1] {
			ans = append(ans, line1[i-1])
			i--
			j--
		} else if ar[i-1][j] >= ar[i][j-1] {
			i--
		} else {
			j--
		}
	}

	for idx := int32(len(ans)) - 1; idx >= 0; idx-- {
		writeInt32(writer, ans[idx])
		if idx > 0 {
			writer.WriteByte(' ')
		}
	}
	writer.WriteByte('\n')
}
