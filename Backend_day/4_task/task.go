package main

import (
	"bufio"
	"fmt"
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
	res = res * negative
	return res
}

type Res struct {
	start  int
	finish int
	value  int
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	N := readInt(reader)
	var line string

	fmt.Fscan(reader, &line)

	_ = N

	res := make(map[int]Res, 27)
	_ = res
	uniq := make(map[rune]bool, 27)
	uniqRes := 0
	for _, v := range line {
		uniq[v] = true
	}
	for key := range uniq {
		_ = key
		uniqRes++
	}

	for D := 1; D <= uniqRes; D++ {
		var counts [26]int
		uniqNow := 0
		left := 0
		right := 0
		ml := 0
		mr := 0
		for right < len(line) {
			value := line[right]
			if counts[value-'a'] == 0 {
				uniqNow++
			}
			counts[value-'a']++
			right++

			for uniqNow > D {
				val := line[left]
				counts[val-'a']--
				if counts[val-'a'] == 0 {
					uniqNow--
				}
				left++
			}

			if right-left > mr-ml {
				mr = right
				ml = left
			}

		}

		r := Res{ml, mr, mr - ml}
		res[D] = r
	}
	maxL := res[1].start
	maxR := res[1].finish
	bestValue := res[1].value
	bestD := 1
	for key, v := range res {

		if v.value*bestD > bestValue*key {
			maxL = v.start
			maxR = v.finish
			bestD = key
			bestValue = v.value
		}
	}
	fmt.Printf("%v %v", maxL+1, maxR)
}
