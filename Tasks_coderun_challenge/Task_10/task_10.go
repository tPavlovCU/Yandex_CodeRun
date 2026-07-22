package main

import (
	"bufio"
	"io"
	"os"
)

type FastScanner struct {
	b []byte
	p int
	n int
}

func NewFastScanner() *FastScanner {
	b, _ := io.ReadAll(os.Stdin)
	return &FastScanner{b: b, n: len(b)}
}

func (f *FastScanner) NextInt() int {
	for f.p < f.n && (f.b[f.p] < '0' || f.b[f.p] > '9') {
		f.p++
	}
	x := 0
	for f.p < f.n && f.b[f.p] >= '0' && f.b[f.p] <= '9' {
		x = x*10 + int(f.b[f.p]-'0')
		f.p++
	}
	return x
}

const MAXNODE = 30000
const MAXCOL = 1100

type DLX struct {
	L [MAXNODE]int
	R [MAXNODE]int
	U [MAXNODE]int
	D [MAXNODE]int
	C [MAXNODE]int

	row [MAXNODE]int

	S [MAXCOL]int

	size int

	ans [300]int
	ac  int

	rowR [5000]int
	rowC [5000]int
	rowV [5000]int

	grid [256]uint8
}

func (d *DLX) init(cols int) {
	d.size = cols

	for i := 0; i <= cols; i++ {
		d.L[i] = i - 1
		d.R[i] = i + 1
		d.U[i] = i
		d.D[i] = i
		d.C[i] = i
		d.S[i] = 0
	}

	d.L[0] = cols
	d.R[cols] = 0
}

func (d *DLX) addRow(id int, cols []int) {
	first := -1

	for _, c := range cols {
		d.size++
		x := d.size

		d.C[x] = c
		d.row[x] = id

		d.S[c]++

		d.D[x] = c
		d.U[x] = d.U[c]
		d.D[d.U[c]] = x
		d.U[c] = x

		if first == -1 {
			first = x
			d.L[x] = x
			d.R[x] = x
		} else {
			d.R[x] = first
			d.L[x] = d.L[first]
			d.R[d.L[first]] = x
			d.L[first] = x
		}
	}
}

func (d *DLX) cover(c int) {
	d.R[d.L[c]] = d.R[c]
	d.L[d.R[c]] = d.L[c]

	for i := d.D[c]; i != c; i = d.D[i] {
		for j := d.R[i]; j != i; j = d.R[j] {
			d.D[d.U[j]] = d.D[j]
			d.U[d.D[j]] = d.U[j]
			d.S[d.C[j]]--
		}
	}
}

func (d *DLX) uncover(c int) {
	for i := d.U[c]; i != c; i = d.U[i] {
		for j := d.L[i]; j != i; j = d.L[j] {
			d.S[d.C[j]]++
			d.D[d.U[j]] = j
			d.U[d.D[j]] = j
		}
	}

	d.R[d.L[c]] = c
	d.L[d.R[c]] = c
}

func (d *DLX) solve(k int) bool {
	if d.R[0] == 0 {
		d.ac = k
		return true
	}

	c := d.R[0]
	best := d.S[c]

	for j := d.R[c]; j != 0; j = d.R[j] {
		if d.S[j] < best {
			best = d.S[j]
			c = j
		}
	}

	if best == 0 {
		return false
	}

	d.cover(c)

	for r := d.D[c]; r != c; r = d.D[r] {
		d.ans[k] = d.row[r]

		for j := d.R[r]; j != r; j = d.R[j] {
			d.cover(d.C[j])
		}

		if d.solve(k + 1) {
			return true
		}

		for j := d.L[r]; j != r; j = d.L[j] {
			d.uncover(d.C[j])
		}
	}

	d.uncover(c)

	return false
}

func writeInt(w *bufio.Writer, x int) {
	if x == 0 {
		w.WriteByte('0')
		return
	}
	var b [12]byte
	p := 12
	for x > 0 {
		p--
		b[p] = byte('0' + x%10)
		x /= 10
	}
	w.Write(b[p:])
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	T := in.NextInt()

	for ; T > 0; T-- {
		n := in.NextInt()
		N := n * n

		var a [256]int

		for i := 0; i < N*N; i++ {
			a[i] = in.NextInt()
		}

		var d DLX

		cols := 4 * N * N
		d.init(cols)

		id := 0

		for r := 0; r < N; r++ {
			for c := 0; c < N; c++ {
				for v := 1; v <= N; v++ {
					if a[r*N+c] != 0 && a[r*N+c] != v {
						continue
					}

					d.rowR[id] = r
					d.rowC[id] = c
					d.rowV[id] = v

					b := (r/n)*n + c/n

					d.addRow(id, []int{
						r*N + c + 1,
						N*N + r*N + v,
						2*N*N + c*N + v,
						3*N*N + b*N + v,
					})

					id++
				}
			}
		}

		if !d.solve(0) {
			out.WriteString("NO\n")
			continue
		}

		for i := 0; i < d.ac; i++ {
			x := d.ans[i]

			r := d.rowR[x]
			c := d.rowC[x]
			v := d.rowV[x]

			d.grid[r*N+c] = uint8(v)
		}

		out.WriteString("YES\n")

		for r := 0; r < N; r++ {
			for c := 0; c < N; c++ {
				if c > 0 {
					out.WriteByte(' ')
				}
				writeInt(out, int(d.grid[r*N+c]))
			}
			out.WriteByte('\n')
		}
	}
}
