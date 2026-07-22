package main

import (
	"bufio"
	"container/heap"
	"os"
	"strconv"
)

type Participant struct {
	Id         int32
	Score      int32
	ProtocolId int32
}

type PriorityQueue []Participant

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].Id == pq[j].Id {
		return pq[i].ProtocolId > pq[j].ProtocolId
	}
	return pq[i].Id < pq[j].Id
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	new := x.(Participant)
	*pq = append(*pq, new)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	deleted := old[n-1]
	*pq = old[0 : n-1]
	return deleted
}

func readInt(reader *bufio.Reader) (int32, bool) {
	var res int32 = 0
	b, err := reader.ReadByte()

	for err == nil && (b == '\n' || b == ' ' || b == '\r') {
		b, err = reader.ReadByte()
	}

	if err != nil {
		return 0, false
	}

	for err == nil && b >= '0' && b <= '9' {
		res = res*10 + int32(b-'0')
		b, err = reader.ReadByte()
	}
	return res, true

}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	M, _ := readInt(reader)

	pq := make(PriorityQueue, 0, M)
	heap.Init(&pq)

	allData := make([]Participant, 0, 5000000)

	startId := make([]int32, M+1)
	endId := make([]int32, M+1)
	cursorId := make([]int32, M+1)
	lastIndx := 0

	for m := int32(0); m < M; m++ {
		L, _ := readInt(reader)
		startId[lastIndx] = int32(len(allData))
		flag := true

		for l := int32(0); l < L; l++ {
			k, _ := readInt(reader)
			v, _ := readInt(reader)
			participant := Participant{Id: k, Score: v, ProtocolId: m}
			allData = append(allData, participant)
			if flag {
				heap.Push(&pq, participant)
				flag = false
			}

		}
		endId[lastIndx] = int32(len(allData)) - 1
		lastIndx += 1
	}

	lastId := int32(-1)
	outBuf := make([]byte, 0, 32)
	for pq.Len() > 0 {
		top := heap.Pop(&pq).(Participant)
		if top.Id > lastId {
			outBuf = outBuf[:0] // Быстрый сброс длины буфера
			outBuf = strconv.AppendInt(outBuf, int64(top.Id), 10)
			outBuf = append(outBuf, ' ')
			outBuf = strconv.AppendInt(outBuf, int64(top.Score), 10)
			outBuf = append(outBuf, '\n')
			writer.Write(outBuf)
			lastId = top.Id
		}
		protocolId := top.ProtocolId
		cursorId[protocolId] += 1
		newItemId := startId[protocolId]
		newItemId += cursorId[protocolId]
		if newItemId <= endId[protocolId] {
			newItem := allData[newItemId]
			heap.Push(&pq, newItem)

		}

	}

}
