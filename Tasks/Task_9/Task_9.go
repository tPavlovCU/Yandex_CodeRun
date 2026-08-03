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
	res = res * negative
	return res
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
func findNeighbours(start int, graph [][]int, visited []bool, result map[int]struct{}) {
	if visited[start] == true {
		return
	} else if len(graph[start]) == 0 {
		result[start] = struct{}{}
		return
	}

	visited[start] = true
	for _, point := range graph[start] {
		result[point] = struct{}{}
		findNeighbours(point, graph, visited, result)
	}
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	N := readInt(reader) //число вершин
	M := readInt(reader) //число ребер

	graph := make([][]int, 0, N+1)
	for n := 0; n < N+1; n++ {
		new := make([]int, 0, N)
		graph = append(graph, new)
	}
	visited := make([]int, N+1)
	visited[0] = 1
	for m := range M {
		_ = m
		p1 := readInt(reader)
		p2 := readInt(reader)
		graph[p1] = append(graph[p1], p2)
		graph[p2] = append(graph[p2], p1)
	}

}
