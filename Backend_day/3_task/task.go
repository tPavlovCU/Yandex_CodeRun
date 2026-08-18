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

func getArray(s string) [26]int {
	var ans [26]int
	for _, value := range s {
		ans[value-'a']++
	}
	return ans
}

type Str struct {
	LeftTail  string
	Middle    [26]int
	RightTail string
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	N := readInt(reader)
	K := readInt(reader)
	strings := make([]string, 0, N)
	_ = K
	for i := 0; i < N; i++ {
		str, _ := reader.ReadString('\n')
		str = str[:len(str)-1]
		strings = append(strings, str)
	}

	result := make(map[Str]bool, N*N+5)
	for startWindow := 0; startWindow < N-K+1; startWindow++ {
		for _, str := range strings {
			strLeftTail := str[:startWindow]
			middle := str[startWindow : startWindow+K]
			strRightTail := str[startWindow+K:]

			s := Str{strLeftTail, getArray(middle), strRightTail}
			_, ok := result[s]

			if ok {
				fmt.Println("YES")
				return
			}
			result[s] = true
		}
	}
	fmt.Println("NO")
}
