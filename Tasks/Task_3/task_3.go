package main

import (
	"bufio"
	"os"
)

func main() {
	writer := bufio.NewWriterSize(os.Stdout, 16*1024)
	defer writer.Flush()

	rBuf := make([]byte, 128*1024)
	var pos, size int

	readNext := func() (int32, bool) {
		for {
			if pos >= size {
				n, err := os.Stdin.Read(rBuf)
				if n == 0 || err != nil {
					return 0, false
				}
				size = n
				pos = 0
			}
			b := rBuf[pos]
			if b >= '0' && b <= '9' {
				break
			}
			pos++
		}

		var res int32 = 0
		for {
			if pos >= size {
				n, _ := os.Stdin.Read(rBuf)
				if n == 0 {
					return res, true
				}
				size = n
				pos = 0
			}
			b := rBuf[pos]
			if b < '0' || b > '9' {
				pos++
				break
			}
			res = res*10 + int32(b-'0')
			pos++
		}
		return res, true
	}

	N_int, ok := readNext()
	if !ok {
		return
	}
	N := int32(N_int)

	arrived := make([]uint64, (N>>6)+2)

	var idx int32 = 0
	var result int32 = 0
	var streak int32 = 0
	N1 := max(N/20, 25)
	for n := int32(0); n < N; n++ {
		newGroup, _ := readNext()

		gIdx := newGroup - 1
		arrived[gIdx>>6] |= uint64(1) << (gIdx & 63)

		for (arrived[idx>>6] & (uint64(1) << (idx & 63))) != 0 {
			idx++
		}

		currentBackstage := n - idx + 1
		if currentBackstage > result {
			result = currentBackstage
			streak = 0
		} else {
			streak += 1

		}

		if currentBackstage+(N-n-1) <= result || streak == N1 {
			break
		}
	}

	if result == 0 {
		writer.WriteByte('0')
	} else {
		var outBuf [16]byte
		p := len(outBuf)
		for result > 0 {
			p--
			outBuf[p] = byte('0' + result%10)
			result /= 10
		}
		writer.Write(outBuf[p:])
	}
}
