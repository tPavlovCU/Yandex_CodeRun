package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type w struct {
	Line1 string
	Line2 string
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	n, _ := reader.ReadString('\n')

	n = strings.TrimSpace(n)
	N, _ := strconv.Atoi(n)

	result := make(map[w]int, N)
	unique := make(map[string]bool, N)

	for i := 0; i < N; i++ {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		for m := 0; m < len(line)-3; m++ {
			lastLine := line[m : m+3]
			nextLine := line[m+1 : m+4]
			newW := w{lastLine, nextLine}
			unique[lastLine] = true
			unique[nextLine] = true

			_, ok := result[newW]

			if !ok {
				result[newW] = 1
			} else {
				result[newW]++
			}
		}
	}
	fmt.Println(len(unique))
	fmt.Println(len(result))

	for key, value := range result {
		fmt.Printf("%v %v %v\n", key.Line1, key.Line2, value)
	}
}
