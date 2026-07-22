package main

import (
	"bufio"
	"os"
	"slices"
)

type Query struct {
	id     int32
	r1, c1 int32
	r2, c2 int32
	x      int32
	ans    int32
}

type FastScanner struct {
	reader *bufio.Reader
	buf    []byte
	pos    int
	size   int
}

func NewFastScanner(r *bufio.Reader) *FastScanner {
	return &FastScanner{
		reader: r,
		buf:    make([]byte, 128*1024),
	}
}

func (fs *FastScanner) NextInt32() (int32, bool) {
	for {
		if fs.pos >= fs.size {
			n, err := fs.reader.Read(fs.buf)
			if n == 0 || err != nil {
				return 0, false
			}
			fs.size = n
			fs.pos = 0
		}
		b := fs.buf[fs.pos]
		if b >= '0' && b <= '9' {
			break
		}
		fs.pos++
	}
	var res int32 = 0
	for {
		if fs.pos >= fs.size {
			n, _ := fs.reader.Read(fs.buf)
			if n == 0 {
				return res, true
			}
			fs.size = n
			fs.pos = 0
		}
		b := fs.buf[fs.pos]
		if b < '0' || b > '9' {
			fs.pos++
			break
		}
		res = res*10 + int32(b-'0')
		fs.pos++
	}
	return res, true
}

func writeInt32(w *bufio.Writer, n int32, buf []byte) {
	if n == 0 {
		w.WriteByte('0')
		return
	}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	w.Write(buf[pos:])
}

var tree []int32

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	fs := NewFastScanner(reader)
	T_val, ok := fs.NextInt32()
	if !ok {
		return
	}
	T := int(T_val)

	cellR := make([]int32, 1000005)
	cellC := make([]int32, 1000005)
	cellVal := make([]int32, 1000005)
	cellOrder := make([]int32, 1000005)

	queries := make([]Query, 0, 250005)
	qOrder := make([]int32, 0, 250005)
	numBuf := make([]byte, 16)

	for t := 0; t < T; t++ {
		n, _ := fs.NextInt32()
		m, _ := fs.NextInt32()
		q, _ := fs.NextInt32()

		cellIdx := 0
		for i := int32(1); i <= n; i++ {
			for j := int32(1); j <= m; j++ {
				v, _ := fs.NextInt32()
				cellR[cellIdx] = i
				cellC[cellIdx] = j
				cellVal[cellIdx] = v
				cellOrder[cellIdx] = int32(cellIdx)
				cellIdx++
			}
		}
		totalCells := cellIdx

		queries = queries[:0]
		qOrder = qOrder[:0]
		for i := int32(0); i < q; i++ {
			r1, _ := fs.NextInt32()
			c1, _ := fs.NextInt32()
			r2, _ := fs.NextInt32()
			c2, _ := fs.NextInt32()
			x, _ := fs.NextInt32()
			queries = append(queries, Query{id: i, r1: r1, c1: c1, r2: r2, c2: c2, x: x})
			qOrder = append(qOrder, i)
		}

		slices.SortFunc(cellOrder[:totalCells], func(i, j int32) int {
			return int(cellVal[j] - cellVal[i])
		})

		slices.SortFunc(qOrder, func(i, j int32) int {
			return int(queries[j].x - queries[i].x)
		})

		colsCount := m + 1
		size := int((n + 1) * colsCount)
		if len(tree) < size {
			tree = make([]int32, size*2)
		} else {
			for i := 0; i < size; i++ {
				tree[i] = 0
			}
		}

		currCellPtr := 0

		query2D := func(r, c int32) int32 {
			var sum int32 = 0
			for i := r; i > 0; i -= i & -i {
				rowOffset := int(i) * int(colsCount)
				for j := c; j > 0; j -= j & -j {
					sum += tree[rowOffset+int(j)]
				}
			}
			return sum
		}

		for _, qId := range qOrder {
			qPtr := &queries[qId]

			for currCellPtr < totalCells {
				cIdx := cellOrder[currCellPtr]
				if cellVal[cIdx] < qPtr.x {
					break
				}

				cr := cellR[cIdx]
				cc := cellC[cIdx]

				for i := cr; i <= n; i += i & -i {
					rowOffset := int(i) * int(colsCount)
					for j := cc; j <= m; j += j & -j {
						tree[rowOffset+int(j)] += 1
					}
				}
				currCellPtr++
			}

			qPtr.ans = query2D(qPtr.r2, qPtr.c2) - query2D(qPtr.r1-1, qPtr.c2) - query2D(qPtr.r2, qPtr.c1-1) + query2D(qPtr.r1-1, qPtr.c1-1)
		}

		for i := int32(0); i < q; i++ {
			writeInt32(writer, queries[i].ans, numBuf)
			writer.WriteByte('\n')
		}
	}
}
