package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)

	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	chars, _ := reader.ReadString('\n')
	chars = strings.TrimSpace(chars)

	needChars := len(chars)
	nowChars := 0

	charsMap := make(map[byte]int, 0)

	for _, value := range chars {
		charsMap[byte(value)] = 0
	}

	badChars := 0
	left := 0
	right := 0
	res := len(line) + 1
	for right < len(line) {
		v, okR := charsMap[line[right]]

		if okR {
			if v == 0 {
				nowChars++
			}
			charsMap[line[right]]++
		} else {
			badChars++
		}

		for nowChars == needChars {
			if right-left+1 < res && badChars == 0 {
				res = right - left + 1
			}

			_, okL := charsMap[line[left]]

			if okL {
				charsMap[line[left]]--
				if charsMap[line[left]] == 0 {
					nowChars--
				}
			} else {
				badChars--
			}
			left++
		}
		right++
	}
	if res == len(line)+1 {
		fmt.Println(0)
		return
	}
	fmt.Println(res)
}
