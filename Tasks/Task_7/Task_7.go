package main

import (
	"bufio"
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

type Edge struct {
	p1 int32
	p2 int32
}

func findNeighbours(value int32, edges []Edge, edgesInverted []Edge, points []int32, pointsInverted []int32, visited []bool, result map[int32]struct{}) {
	visitedCopy := make([]bool, 0, len(visited))
	visitedCopy = append(visitedCopy, visited...)

	result[value] = struct{}{}
	if points[value] == -1 || visited[value] == true {

	} else {
		lenn := int32(len(edges))
		result[value] = struct{}{}
		for start := points[value]; start < lenn; start++ {
			if edges[start].p1 != value {
				break
			}
			el := edges[start].p2
			visited[value] = true
			findNeighbours(el, edges, edgesInverted, points, pointsInverted, visited, result)
		}

	}
	if pointsInverted[value] == -1 || visitedCopy[value] == true {

	} else {
		lenn := int32(len(edges))
		result[value] = struct{}{}
		for start := pointsInverted[value]; start < lenn; start++ {
			if edgesInverted[start].p1 != value {
				break
			}
			el := edgesInverted[start].p2
			visitedCopy[value] = true
			findNeighbours(el, edges, edgesInverted, points, pointsInverted, visited, result)
		}

	}
	return
}
func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<11)
	writer := bufio.NewWriterSize(os.Stdout, 1<<4)
	defer writer.Flush()
	_ = reader

	N := readInt(reader)
	_ = N
	M := readInt(reader)

	edges := make([]Edge, 0, M)
	edgesInverted := make([]Edge, 0, M)
	points := make([]int32, 0, N+2)
	pointsInverted := make([]int32, 0, N+2)
	for n := int32(0); n < N+1; n++ {
		points = append(points, -1)
		pointsInverted = append(pointsInverted, -1)
	}

	for m := int32(0); m < M; m++ {
		p1 := readInt(reader)
		p2 := readInt(reader)
		edge := Edge{p1, p2}
		edges = append(edges, edge)
		edgeInverted := Edge{p2, p1}
		edgesInverted = append(edgesInverted, edgeInverted)
	}

	sort.Slice(edges, func(i, j int) bool {
		edgeI := edges[i]
		edgeJ := edges[j]
		if edgeI.p1 < edgeJ.p1 {
			return true
		} else if edgeI.p1 == edgeJ.p1 {
			return edgeI.p2 < edgeJ.p2
		} else {
			return false
		}

	})
	sort.Slice(edgesInverted, func(i, j int) bool {
		edgeI := edgesInverted[i]
		edgeJ := edgesInverted[j]
		if edgeI.p1 < edgeJ.p1 {
			return true
		} else if edgeI.p1 == edgeJ.p1 {
			return edgeI.p2 < edgeJ.p2
		} else {
			return false
		}

	})

	lastValue := int32(0)
	for idx, edge := range edges {
		p1 := edge.p1

		if p1 > lastValue {
			points[p1] = int32(idx)
			lastValue = p1
		}
	}
	lastValue = int32(0)
	for idx, edge := range edgesInverted {
		p1 := edge.p1

		if p1 > lastValue {
			pointsInverted[p1] = int32(idx)
			lastValue = p1
		}
	}

	visited := make([]bool, N+1)
	result := make(map[int32]struct{}, N+1)
	findNeighbours(1, edges, edgesInverted, points, pointsInverted, visited, result)

	ar := make([]int, 0, len(result))
	for key := range result {
		ar = append(ar, int(key))
	}
	sort.Ints(ar)
	writeInt32(writer, int32(len(result)))
	writer.WriteByte('\n')
	for idx, value := range ar {
		if idx != len(ar)-1 {
			writeInt32(writer, int32(value))
			writer.WriteByte(' ')
		} else {
			writeInt32(writer, int32(value))
		}
	}
}
