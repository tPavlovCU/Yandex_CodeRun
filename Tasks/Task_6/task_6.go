package main

import (
	"bufio"
	"io"
	"os"
)

type FastScanner struct {
	reader io.Reader
	buf    []byte
	pos    int
	size   int
}

func NewFastScanner(r io.Reader) *FastScanner {
	return &FastScanner{
		reader: r,
		buf:    make([]byte, 128*1024),
	}
}

func (fs *FastScanner) NextInt64() (int64, bool) {
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
	var res int64 = 0
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
		res = res*10 + int64(b-'0')
		fs.pos++
	}
	return res, true
}

func writeInt64(w *bufio.Writer, n int64, buf []byte) {
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

var a []int64
var bArray []int64
var deque []int32

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	fs := NewFastScanner(reader)
	tVal, ok := fs.NextInt64()
	if !ok {
		return
	}
	T := int(tVal)

	a = make([]int64, 500005)
	bArray = make([]int64, 500005)
	deque = make([]int32, 500005)
	numBuf := make([]byte, 24)

	for t := 0; t < T; t++ {
		nVal, _ := fs.NextInt64()
		kVal, _ := fs.NextInt64()
		n := int32(nVal)
		k := int32(kVal)

		for i := int32(0); i < n; i++ {
			a[i], _ = fs.NextInt64()
		}

		var curSum int64 = 0
		var curXor int64 = 0
		for i := int32(0); i < k; i++ {
			curSum += a[i]
			curXor ^= a[i]
		}
		bArray[0] = curSum - curXor

		for l := int32(1); l <= n-k; l++ {
			curSum = curSum - a[l-1] + a[l+k-1]
			curXor = curXor ^ a[l-1] ^ a[l+k-1]
			bArray[l] = curSum - curXor
		}

		headPtr := 0
		tailPtr := 0
		nextL := int32(0)
		limitL := n - k

		for i := int32(0); i < n; i++ {
			lStart := i - k + 1
			if lStart < 0 {
				lStart = 0
			}
			lEnd := i
			if lEnd > limitL {
				lEnd = limitL
			}

			for nextL <= lEnd {
				val := bArray[nextL]
				for tailPtr > headPtr && bArray[deque[tailPtr-1]] >= val {
					tailPtr--
				}
				deque[tailPtr] = nextL
				tailPtr++
				nextL++
			}

			for headPtr < tailPtr && deque[headPtr] < lStart {
				headPtr++
			}

			if i > 0 {
				writer.WriteByte(' ')
			}
			writeInt64(writer, bArray[deque[headPtr]], numBuf)
		}
		writer.WriteByte('\n')
	}
}
