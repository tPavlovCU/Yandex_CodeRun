package main

import (
	"bufio"
	"math/bits"
	"os"
)

type Matrix [2][2]uint64

var powers [63]Matrix

func multiplyMatrix(A, B Matrix, M uint64) Matrix {
	return Matrix{
		{(A[0][0]*B[0][0] + A[0][1]*B[1][0]) % M, (A[0][0]*B[0][1] + A[0][1]*B[1][1]) % M},
		{(A[1][0]*B[0][0] + A[1][1]*B[1][0]) % M, (A[1][0]*B[0][1] + A[1][1]*B[1][1]) % M},
	}
}

func readUint(reader *bufio.Reader) (uint64, bool) {
	var res uint64 = 0
	b, err := reader.ReadByte()
	for err == nil && (b == '\n' || b == ' ' || b == '\r') {
		b, err = reader.ReadByte()
	}
	if err != nil {
		return 0, false
	}
	for err == nil && b >= '0' && b <= '9' {
		res = res*10 + uint64(b-'0')
		b, err = reader.ReadByte()
	}
	return res, true
}

func writeUint64(w *bufio.Writer, n uint64) {
	if n == 0 {
		w.WriteByte('0')
		return
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	w.Write(buf[pos:])
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	a, _ := readUint(reader)
	b, _ := readUint(reader)
	M, _ := readUint(reader)
	Q, _ := readUint(reader)

	powers[0] = Matrix{
		{a % M, b % M},
		{1, 0},
	}

	for i := 1; i < 63; i++ {
		powers[i] = multiplyMatrix(powers[i-1], powers[i-1], M)
	}

	for q := uint64(0); q < Q; q++ {
		x, _ := readUint(reader)
		y, _ := readUint(reader)
		pow, _ := readUint(reader)

		if pow == 0 {
			writeUint64(writer, x%M)
			w := writer
			w.WriteByte('\n')
			continue
		}
		if pow == 1 {
			writeUint64(writer, y%M)
			w := writer
			w.WriteByte('\n')
			continue
		}

		pow--

		v0 := y % M
		v1 := x % M

		var bit int

		for pow > 0 {
			trailing := bits.TrailingZeros64(pow)
			pow >>= trailing
			bit += trailing

			mat := &powers[bit]
			nextV0 := (mat[0][0]*v0 + mat[0][1]*v1) % M
			nextV1 := (mat[1][0]*v0 + mat[1][1]*v1) % M

			v0 = nextV0
			v1 = nextV1

			pow >>= 1
			bit++
		}

		writeUint64(writer, v0)
		w := writer
		w.WriteByte('\n')
	}
}
