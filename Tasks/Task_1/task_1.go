package main

import (
	"bufio"
	"os"
	"slices"
)

func readInt(reader *bufio.Reader) (int64, bool) {
	var res int64 = 0
	b, err := reader.ReadByte()

	for err == nil && (b == '\n' || b == ' ' || b == '\r') {
		b, err = reader.ReadByte()
	}

	if err != nil {
		return 0, false
	}

	for err == nil && b >= '0' && b <= '9' {
		res = res*10 + int64(b-'0')
		b, err = reader.ReadByte()
	}
	return res, true
}

func writeInt64(w *bufio.Writer, n int64) {
	if n == 0 {
		w.WriteByte('0')
		return
	}
	var buf [24]byte
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

	N, _ := readInt(reader)
	T, _ := readInt(reader)

	if N > T/2 {
		timeline := make([]int64, T+1)

		for m := int64(0); m < N; m++ {
			A, _ := readInt(reader)
			F, _ := readInt(reader)
			S, _ := readInt(reader)

			timeline[A] += S
			timeline[F] -= S
		}

		var n int64 = 0
		var maximum int64 = 0
		for _, now := range timeline {
			n += now
			if n > maximum {
				maximum = n
			}
		}
		writeInt64(writer, maximum)
	} else {

		timelineMap := make(map[int64]int64, int(N*2))

		for m := int64(0); m < N; m++ {
			A, _ := readInt(reader)
			F, _ := readInt(reader)
			S, _ := readInt(reader)

			timelineMap[A] += S
			timelineMap[F] -= S
		}

		keys := make([]int64, 0, int(N*2))
		for key := range timelineMap {
			keys = append(keys, key)
		}

		slices.Sort(keys)

		var now int64 = 0
		var maximum int64 = 0
		for _, value := range keys {
			now += timelineMap[value]
			if now > maximum {
				maximum = now
			}
		}
		writeInt64(writer, maximum)
	}
}
