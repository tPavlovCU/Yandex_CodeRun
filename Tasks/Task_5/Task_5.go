package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	N := readInt(reader)
	ar := make([][]int32, 0, N+1)

	for n := int32(0); n < N+1; n++ {
		new := make([]int32, 0, N+1)
		for m := int32(0); m < N+1; m++ {
			new = append(new, 30001)
		}
		ar = append(ar, new)
	}
	ar[0][0] = 0
	prices := make([]int32, 0, N+1)
	for days := int32(0); days < N; days++ {
		item := readInt(reader)
		prices = append(prices, item)
		for talons := int32(0); talons < N; talons++ {
			now := ar[days][talons]
			if talons > days {
				ar[days][talons] = -1
			} else {
				if item > 100 {
					ar[days+1][talons+1] = min(now+item, ar[days+1][talons+1])
				} else {
					ar[days+1][talons] = min(now+item, ar[days+1][talons])
				}
				if talons > 0 {
					ar[days+1][talons-1] = min(now, ar[days+1][talons-1])
				}

			}

		}
	}

	minsum := int32(30001)
	talons := int32(0)
	for m := int32(0); m < N; m++ {
		if ar[N][m] <= minsum {
			minsum = ar[N][m]
			talons = m
		}
	}
	writeInt32(writer, minsum)
	writer.WriteByte('\n')

	used := make([]int, 0, N+1)
	n := talons
	fmt.Println(prices)
	for m := N; m > 0; m-- {
		if ar[m][n] == ar[m-1][n+1] {
			used = append(used, int(m))
			n++
		} else if prices[m-1] > 100 {
			n--
		} else {
			n = n
		}
	}
	sort.Ints(used)
	writeInt32(writer, talons)
	writer.WriteByte(' ')
	writeInt32(writer, int32(len(used)))
	writer.WriteByte('\n')
	for idx, value := range used {
		if idx == len(used)-1 {
			writeInt32(writer, int32(value))
		} else {
			writeInt32(writer, int32(value))
			writer.WriteByte(' ')
		}
	}

}
