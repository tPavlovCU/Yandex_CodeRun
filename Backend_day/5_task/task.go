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

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)

	N := readInt(reader)

	mins := make([]int, N)
	maxs := make([]int, N)

	nums := make([]int, 0, N)

	for m := range N {
		_ = m
		newItem := readInt(reader)
		nums = append(nums, newItem)
	}
	mins[0] = nums[0]
	for idx, value := range nums {
		if idx > 0 {
			mins[idx] = min(mins[idx-1], value)
		}
	}

	maxs[len(maxs)-1] = nums[len(nums)-1]
	for i := len(maxs) - 2; i > -1; i-- {
		maxs[i] = max(maxs[i+1], nums[i])
	}

	k := 0
	for i := 1; i < N; i++ {
		if mins[i-1] >= maxs[i] {
			k = i
		}
	}
	fmt.Println(k)

}
