package main

import (
	"bufio"
	"os"
	"sort"
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

func writeInt(w *bufio.Writer, n int) {
	if n == 0 {
		w.WriteByte('0')
		return
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	w.Write(buf[pos:])
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func manhattan(x1, y1, x2, y2 int64) int64 {
	return abs(x1-x2) + abs(y1-y2)
}

type Edge struct {
	cost int64
	u, v int
}

type DSU struct {
	parent []int
	rank   []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	return &DSU{parent, rank}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) bool {
	px, py := d.Find(x), d.Find(y)
	if px == py {
		return false
	}
	if d.rank[px] < d.rank[py] {
		px, py = py, px
	}
	d.parent[py] = px
	if d.rank[px] == d.rank[py] {
		d.rank[px]++
	}
	return true
}

type Pair struct {
	a, b int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	t, _ := readInt(reader)

	for tc := 0; tc < int(t); tc++ {
		n, _ := readInt(reader)

		points := make([][3]int64, n)
		for i := int64(0); i < n; i++ {
			points[i][0], _ = readInt(reader)
			points[i][1], _ = readInt(reader)
			points[i][2], _ = readInt(reader)
		}

		edgeMap := make(map[Pair]int64)

		for i := int64(0); i < n; i++ {
			x1, y1, r1 := points[i][0], points[i][1], points[i][2]
			for j := int64(0); j < n; j++ {
				if i == j {
					continue
				}
				x2, y2 := points[j][0], points[j][1]
				d := manhattan(x1, y1, x2, y2)
				if d <= r1 {
					u, v := i, j
					if u > v {
						u, v = v, u
					}
					key := Pair{int(u), int(v)}
					if prev, exists := edgeMap[key]; !exists || d < prev {
						edgeMap[key] = d
					}
				}
			}
		}

		edges := make([]Edge, 0, len(edgeMap))
		for key, cost := range edgeMap {
			edges = append(edges, Edge{cost, key.a, key.b})
		}

		sort.Slice(edges, func(i, j int) bool {
			return edges[i].cost < edges[j].cost
		})

		dsu := NewDSU(int(n))
		var totalCost int64
		var resultEdges []Edge

		for _, e := range edges {
			if dsu.Union(e.u, e.v) {
				totalCost += e.cost
				resultEdges = append(resultEdges, e)
			}
		}

		writeInt64(writer, totalCost)
		writer.WriteByte('\n')
		writeInt(writer, len(resultEdges))
		writer.WriteByte('\n')
		for _, e := range resultEdges {
			writeInt(writer, e.u+1)
			writer.WriteByte(' ')
			writeInt(writer, e.v+1)
			writer.WriteByte('\n')
		}
	}
}
