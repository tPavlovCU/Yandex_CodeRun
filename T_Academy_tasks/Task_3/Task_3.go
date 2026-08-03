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

	return res * negative
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
		n = -n
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	n := readInt(reader)
	k := readInt(reader)

	var sBytes []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		if b != '\n' && b != '\r' && b != ' ' {
			sBytes = append(sBytes, b)

			for len(sBytes) < n {
				b, _ = reader.ReadByte()
				sBytes = append(sBytes, b)
			}
			break
		}
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = readInt(reader)
	}

	dp := make([]int, n*n*(k+1))

	get := func(l, r, x int) int {
		return dp[(l*n+r)*(k+1)+x]
	}

	set := func(l, r, x, val int) {
		dp[(l*n+r)*(k+1)+x] = val
	}

	for i := 0; i < n; i++ {
		val := a[i]
		if val < 0 {
			val = 0
		}

		for x := 0; x <= k; x++ {
			set(i, i, x, val)
		}
	}

	for length := 2; length <= n; length++ {
		for l := 0; l+length <= n; l++ {
			r := l + length - 1

			for x := 0; x <= k; x++ {

				best := 0

				best = max(best, get(l+1, r, x))

				best = max(best, get(l, r-1, x))

				if sBytes[l] == sBytes[r] {

					val := a[l] + a[r]

					if l+1 <= r-1 {
						val += get(l+1, r-1, x)
					}

					best = max(best, val)

				} else if x > 0 {

					val := a[l] + a[r]

					if l+1 <= r-1 {
						val += get(l+1, r-1, x-1)
					}

					best = max(best, val)
				}

				set(l, r, x, best)
			}
		}
	}

	writeInt(writer, get(0, n-1, k))
}
