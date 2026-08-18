package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	Point   int
	To      int
	visited bool
}

func readInt(reader *bufio.Reader) int {
	b, err := reader.ReadByte()
	res := 0
	if err != nil {
		fmt.Println("error", err)
		return 0
	}

	for b < '0' || b > '9' {
		b, _ = reader.ReadByte()
	}

	for b <= '9' && b >= '0' {
		res = res*10 + int(b) - '0'
		b, _ = reader.ReadByte()
	}
	return res
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)

	n := readInt(reader)
	m := readInt(reader)
	allPoints := make(map[int][]Edge, 2*m)
	uniqPoints := make([]int, 0, n)
	for m := range n {
		uniqPoints = append(uniqPoints, m+1)
	}

	for range m {
		p1 := readInt(reader)
		p2 := readInt(reader)
		edge1 := Edge{p1, p2, false}
		edge2 := Edge{p2, p1, false}

		edgesP1 := allPoints[p1]
		edgesP2 := allPoints[p2]

		edgesP1 = append(edgesP1, edge1)
		edgesP2 = append(edgesP2, edge2)

		allPoints[p1] = edgesP1
		allPoints[p2] = edgesP2
	}

	group1 := make(map[int]struct{}, n)
	group2 := make(map[int]struct{}, n)

	for _, Point := range uniqPoints {
		edges := allPoints[Point]
		group := 1
		otherGroup := 1
		if _, ok := group1[Point]; ok {
			group = 1
			otherGroup = 1
			group1[Point] = struct{}{}
		} else if _, ok := group2[Point]; ok {
			group = 2
			otherGroup = 2
			group2[Point] = struct{}{}
		} else {
			group = 1
			otherGroup = 1
			group1[Point] = struct{}{}
		}

		for _, PointData := range edges {
			secondPoint := PointData.To
			if group == 1 {
				if _, ok := group2[secondPoint]; ok {
					fmt.Println(-1)
					return
				}
			} else {
				if _, ok := group1[secondPoint]; ok {
					fmt.Println(-1)
					return
				}
			}

			if otherGroup == 1 {
				group1[secondPoint] = struct{}{}
			} else {
				group2[secondPoint] = struct{}{}
			}
		}
	}

	firstGroup := make([]int, 0, len(group1))
	secondGroup := make([]int, 0, len(group2))

	for key := range group1 {
		firstGroup = append(firstGroup, key)
	}

	for key := range group2 {
		secondGroup = append(secondGroup, key)
	}
	if len(group1) == 0 {
		fmt.Println(-1)
		return
	}
	sort.Ints(firstGroup)
	sort.Ints(secondGroup)
	fmt.Println(len(group1))
	for _, v := range firstGroup {
		fmt.Printf("%v ", v)
	}
	fmt.Printf("\n")
	for _, v := range secondGroup {
		fmt.Printf("%v ", v)
	}
}
