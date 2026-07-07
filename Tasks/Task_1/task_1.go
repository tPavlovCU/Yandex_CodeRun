package main

import (
	"bufio"
	"os"
	"sort"
	"strconv"
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

func main() {
	/*
	  Пример ввода и вывода числа n, где -10^9 < n < 10^9:
	*/

	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	N, _ := readInt(reader)
	T, _ := readInt(reader)

	if N > T/2 {

		timeline := make([]int64, T+1, T+1)

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

		writer.WriteString(strconv.Itoa(int(maximum)))
	} else {
		timelineMap := make(map[int64]int64, N*2)

		for m := int64(0); m < N; m++ {
			A, _ := readInt(reader)
			F, _ := readInt(reader)
			S, _ := readInt(reader)

			timelineMap[A] += S
			timelineMap[F] -= S
		}

		keys := make([]int64, 0, N*2)

		var now int64 = 0
		var maximum int64 = 0

		for key := range timelineMap {
			keys = append(keys, key)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, value := range keys {

			now += timelineMap[value]
			if now > maximum {
				maximum = now
			}
		}
		writer.WriteString(strconv.Itoa(int(maximum)))
	}
}
