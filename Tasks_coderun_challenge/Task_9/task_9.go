package main

import (
	"bufio"
	"io"
	"os"
	"sort"
)

type FastScanner struct {
	data []byte
	pos  int
	n    int
}

func NewFastScanner() *FastScanner {
	data, _ := io.ReadAll(os.Stdin)
	return &FastScanner{
		data: data,
		n:    len(data),
	}
}

func (fs *FastScanner) NextInt() int32 {
	for fs.pos < fs.n {
		c := fs.data[fs.pos]
		if c >= '0' && c <= '9' {
			break
		}
		fs.pos++
	}

	var x int32
	for fs.pos < fs.n {
		c := fs.data[fs.pos]
		if c < '0' || c > '9' {
			break
		}
		x = x*10 + int32(c-'0')
		fs.pos++
	}

	return x
}

func writeInt(w *bufio.Writer, x int32) {
	if x == 0 {
		w.WriteByte('0')
		w.WriteByte('\n')
		return
	}

	var buf [12]byte
	i := len(buf)

	for x > 0 {
		i--
		buf[i] = byte('0' + x%10)
		x /= 10
	}

	w.Write(buf[i:])
	w.WriteByte('\n')
}

type Interval struct {
	s  int32
	e  int32
	si int
	ei int
}

type Query struct {
	l  int32
	r  int32
	ri int
}

type SegTree struct {
	t []int32
}

func NewSegTree(n int) *SegTree {
	return &SegTree{
		t: make([]int32, n<<2),
	}
}

func (st *SegTree) update(v, l, r, q int, val int32) {
	if q <= l {
		if st.t[v] < val {
			st.t[v] = val
		}
		return
	}

	m := (l + r) >> 1

	if q <= m {
		st.update(v<<1, l, m, q, val)
	} else {
		st.update(v<<1|1, m+1, r, q, val)
	}

	if st.t[v<<1] > st.t[v] {
		st.t[v] = st.t[v<<1]
	}
	if st.t[v<<1|1] > st.t[v] {
		st.t[v] = st.t[v<<1|1]
	}
}

func (st *SegTree) query(v, l, r, q int, cur int32) int32 {
	if st.t[v] > cur {
		cur = st.t[v]
	}

	if l == r {
		return cur
	}

	m := (l + r) >> 1

	if q <= m {
		return st.query(v<<1, l, m, q, cur)
	}

	return st.query(v<<1|1, m+1, r, q, cur)
}

func lowerBound(a []int32, x int32) int {
	l, r := 0, len(a)

	for l < r {
		m := (l + r) >> 1
		if a[m] < x {
			l = m + 1
		} else {
			r = m
		}
	}

	return l
}

func main() {
	fs := NewFastScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	tc := int(fs.NextInt())

	for ; tc > 0; tc-- {
		n := int(fs.NextInt())
		t := fs.NextInt()

		normal := make([]Interval, 0, n)
		wrap := make([]Query, 0)

		coords := make([]int32, 0, n*2+1)

		for i := 0; i < n; i++ {
			s := fs.NextInt()
			d := fs.NextInt()

			e := s + d

			if e <= t {
				normal = append(normal, Interval{
					s: s,
					e: e,
				})
				coords = append(coords, s, e)
			} else {
				wrap = append(wrap, Query{
					l: e - t,
					r: s,
				})
				coords = append(coords, s)
			}
		}

		if len(normal) == 0 {
			if len(wrap) > 0 {
				writeInt(out, 1)
			} else {
				writeInt(out, 0)
			}
			continue
		}

		sort.Slice(coords, func(i, j int) bool {
			return coords[i] < coords[j]
		})

		m := 0
		for _, x := range coords {
			if m == 0 || coords[m-1] != x {
				coords[m] = x
				m++
			}
		}
		coords = coords[:m]

		for i := range normal {
			normal[i].si = lowerBound(coords, normal[i].s)
			normal[i].ei = lowerBound(coords, normal[i].e)
		}

		for i := range wrap {
			wrap[i].ri = lowerBound(coords, wrap[i].r)
		}

		sort.Slice(wrap, func(i, j int) bool {
			return wrap[i].l > wrap[j].l
		})

		st := NewSegTree(len(coords))

		ptr := len(normal) - 1
		var ans int32

		for _, q := range wrap {
			for ptr >= 0 && normal[ptr].s >= q.l {
				best := st.query(
					1,
					0,
					len(coords)-1,
					normal[ptr].si,
					0,
				) + 1

				st.update(
					1,
					0,
					len(coords)-1,
					normal[ptr].ei,
					best,
				)

				ptr--
			}

			cur := st.query(
				1,
				0,
				len(coords)-1,
				q.ri,
				0,
			) + 1

			if cur > ans {
				ans = cur
			}
		}

		st = NewSegTree(len(coords))
		ptr = len(normal) - 1

		for ptr >= 0 {
			best := st.query(
				1,
				0,
				len(coords)-1,
				normal[ptr].si,
				0,
			) + 1

			st.update(
				1,
				0,
				len(coords)-1,
				normal[ptr].ei,
				best,
			)

			ptr--
		}

		cur := st.query(
			1,
			0,
			len(coords)-1,
			len(coords)-1,
			0,
		)

		if cur > ans {
			ans = cur
		}

		writeInt(out, ans)
	}
}
