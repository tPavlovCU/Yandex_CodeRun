package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
)

type TourResponse struct {
	StartIsland int16  `json:"start_island"`
	MaxStamps   int16  `json:"max_stamps"`
	TourId      string `json:"tour_id"`
}

type RoutesResponse struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	From int16 `json:"from"`
	To   int16 `json:"to"`
}

type StampsResponse struct {
	Stamps []Stamp `json:"stamps"`
}

type Stamp struct {
	StampId  int16   `json:"stamp_id"`
	IslandId int16   `json:"island_id"`
	Rarity   int16   `json:"rarity"`
	Requires []int16 `json:"requires"`
}

type Edge struct {
	to   int16
	next int16
	id   int16
}

func readLine(reader *bufio.Reader) (string, bool) {
	var res []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if len(res) > 0 {
				return string(res), true
			}
			return "", false
		}
		if b == '\n' || b == '\r' {
			if len(res) > 0 {
				return string(res), true
			}
			continue
		}
		res = append(res, b)
	}
}

func writeInt16(w *bufio.Writer, n int16) {
	if n == 0 {
		w.WriteByte('0')
		return
	}
	var buf [12]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	w.Write(buf[pos:])
}

func main() {
	reader := bufio.Reader{}
	reader.Reset(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	baseURL, ok := readLine(&reader)
	if !ok {
		writer.WriteByte('0')
		writer.WriteByte('\n')
		return
	}
	tourIDStr, ok := readLine(&reader)
	if !ok {
		writer.WriteByte('0')
		writer.WriteByte('\n')
		return
	}

	client := &http.Client{}

	resp, err := client.Get(baseURL + "/tour/" + tourIDStr)
	if err != nil {
		writer.WriteByte('0')
		writer.WriteByte('\n')
		return
	}
	var tourData TourResponse
	json.NewDecoder(resp.Body).Decode(&tourData)
	resp.Body.Close()

	resp, err = client.Get(baseURL + "/routes?tour_id=" + tourIDStr)
	if err != nil {
		writer.WriteByte('0')
		writer.WriteByte('\n')
		return
	}
	var routesData RoutesResponse
	json.NewDecoder(resp.Body).Decode(&routesData)
	resp.Body.Close()

	resp, err = client.Get(baseURL + "/stamps?tour_id=" + tourIDStr)
	if err != nil {
		writer.WriteByte('0')
		writer.WriteByte('\n')
		return
	}
	var stampsData StampsResponse
	json.NewDecoder(resp.Body).Decode(&stampsData)
	resp.Body.Close()

	var maxIslands int16 = 0
	for i := 0; i < len(routesData.Routes); i++ {
		r := routesData.Routes[i]
		if r.From > maxIslands {
			maxIslands = r.From
		}
		if r.To > maxIslands {
			maxIslands = r.To
		}
	}
	for i := 0; i < len(stampsData.Stamps); i++ {
		s := stampsData.Stamps[i]
		if s.IslandId > maxIslands {
			maxIslands = s.IslandId
		}
	}

	adjHead := make([]int16, maxIslands+1)
	for i := range adjHead {
		adjHead[i] = -1
	}
	adjEdges := make([]Edge, len(routesData.Routes))
	for i := 0; i < len(routesData.Routes); i++ {
		r := routesData.Routes[i]
		adjEdges[i] = Edge{to: r.To, next: adjHead[r.From]}
		adjHead[r.From] = int16(i)
	}

	reachable := make([]bool, maxIslands+1)
	if tourData.StartIsland <= maxIslands {
		q := make([]int16, 0, maxIslands+1)
		q = append(q, tourData.StartIsland)
		reachable[tourData.StartIsland] = true
		head := 0
		for head < len(q) {
			u := q[head]
			head++
			e := adjHead[u]
			for e != -1 {
				v := adjEdges[e].to
				if !reachable[v] {
					reachable[v] = true
					q = append(q, v)
				}
				e = adjEdges[e].next
			}
		}
	}

	var maxStampID int16 = 0
	for i := 0; i < len(stampsData.Stamps); i++ {
		if stampsData.Stamps[i].StampId > maxStampID {
			maxStampID = stampsData.Stamps[i].StampId
		}
	}

	stampIDToIdx := make([]int16, maxStampID+1)
	for i := range stampIDToIdx {
		stampIDToIdx[i] = -1
	}

	var validStamps []Stamp
	for i := 0; i < len(stampsData.Stamps); i++ {
		s := stampsData.Stamps[i]
		if s.IslandId <= maxIslands {
			if reachable[s.IslandId] {
				stampIDToIdx[s.StampId] = int16(len(validStamps))
				validStamps = append(validStamps, s)
			}
		}
	}

	n := int16(len(validStamps))
	k := tourData.MaxStamps
	if n == 0 || k <= 0 {
		writer.WriteByte('0')
		writer.WriteByte('\n')
		return
	}

	stampHead := make([]int16, n+1)
	for i := range stampHead {
		stampHead[i] = -1
	}
	stampEdges := make([]Edge, n*2+2)
	var edgeCount int16 = 0
	inDegree := make([]int16, n+1)

	for i := int16(0); i < n; i++ {
		s := validStamps[i]
		if len(s.Requires) > 0 {
			pID := s.Requires[0]
			if pID <= maxStampID {
				pIdx := stampIDToIdx[pID]
				if pIdx != -1 {
					stampEdges[edgeCount] = Edge{to: i + 1, next: stampHead[pIdx+1], id: edgeCount}
					stampHead[pIdx+1] = edgeCount
					edgeCount++
					inDegree[i+1]++
				} else {
					inDegree[i+1] = -1
				}
			}
		}
	}

	for i := int16(1); i <= n; i++ {
		if inDegree[i] == 0 {
			stampEdges[edgeCount] = Edge{to: i, next: stampHead[0], id: edgeCount}
			stampHead[0] = edgeCount
			edgeCount++
		}
	}

	order := make([]int16, 0, n+1)
	visited := make([]bool, n+1)
	stack := make([]int16, 0, n+1)
	edgePtr := make([]int16, n+1)
	for i := range edgePtr {
		edgePtr[i] = stampHead[i]
	}

	stack = append(stack, 0)
	visited[0] = true

	for len(stack) > 0 {
		u := stack[len(stack)-1]
		e := edgePtr[u]
		if e != -1 {
			edgePtr[u] = stampEdges[e].next
			v := stampEdges[e].to
			if inDegree[v] != -1 {
				if !visited[v] {
					visited[v] = true
					stack = append(stack, v)
				}
			}
		} else {
			order = append(order, u)
			stack = stack[:len(stack)-1]
		}
	}

	stride := k + 1
	dp := make([][]int16, n+1)
	for i := range dp {
		dp[i] = make([]int16, stride)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}

	history := make([][]int16, edgeCount)
	for i := range history {
		history[i] = make([]int16, stride)
		for j := range history[i] {
			history[i][j] = -1
		}
	}

	currentDP := make([]int16, stride)
	nextDP := make([]int16, stride)
	bestChoice := make([]int16, stride)

	for i := 0; i < len(order); i++ {
		u := order[i]

		for j := int16(0); j <= k; j++ {
			currentDP[j] = -1
		}
		currentDP[0] = 0

		e := stampHead[u]
		for e != -1 {
			v := stampEdges[e].to
			edgeID := stampEdges[e].id
			if inDegree[v] == -1 {
				e = stampEdges[e].next
				continue
			}

			for j := int16(0); j <= k; j++ {
				nextDP[j] = -1
				bestChoice[j] = -1
			}

			for j := int16(0); j <= k; j++ {
				if currentDP[j] != -1 {
					for k2 := int16(0); j+k2 <= k; k2++ {
						if dp[v][k2] != -1 {
							score := currentDP[j] + dp[v][k2]
							if score > nextDP[j+k2] {
								nextDP[j+k2] = score
								bestChoice[j+k2] = k2
							}
						}
					}
				}
			}

			for j := int16(0); j <= k; j++ {
				currentDP[j] = nextDP[j]
				if bestChoice[j] != -1 {
					history[edgeID][j] = bestChoice[j]
				}
			}
			e = stampEdges[e].next
		}

		if u == 0 {
			for j := int16(0); j <= k; j++ {
				dp[u][j] = currentDP[j]
			}
		} else {
			weight := validStamps[u-1].Rarity
			for j := int16(1); j <= k; j++ {
				if currentDP[j-1] != -1 {
					dp[u][j] = currentDP[j-1] + weight
				}
			}
			dp[u][0] = 0
		}
	}

	var bestK int16 = 0
	var maxScore int16 = -1
	for j := int16(0); j <= k; j++ {
		if dp[0][j] > maxScore {
			maxScore = dp[0][j]
			bestK = j
		}
	}

	if maxScore <= 0 {
		writer.WriteByte('0')
		writer.WriteByte('\n')
		return
	}

	chosen := make([]bool, n+1)

	var recover func(int16, int16)
	recover = func(u int16, remK int16) {
		if remK <= 0 {
			return
		}
		if u != 0 {
			chosen[u] = true
			remK--
		}

		var edges []Edge
		e := stampHead[u]
		for e != -1 {
			if inDegree[stampEdges[e].to] != -1 {
				edges = append(edges, stampEdges[e])
			}
			e = stampEdges[e].next
		}

		currRem := remK
		for idx := len(edges) - 1; idx >= 0; idx-- {
			edge := edges[idx]
			childK := history[edge.id][currRem]
			if childK != -1 {
				recover(edge.to, childK)
				currRem -= childK
			}
		}
	}

	recover(0, bestK)

	var finalOrder []int16
	var collect func(int16)
	collect = func(u int16) {
		if u != 0 {
			if chosen[u] {
				finalOrder = append(finalOrder, validStamps[u-1].StampId)
			}
		}
		e := stampHead[u]
		var children []int16
		for e != -1 {
			v := stampEdges[e].to
			if inDegree[v] != -1 {
				children = append(children, v)
			}
			e = stampEdges[e].next
		}
		for idx := len(children) - 1; idx >= 0; idx-- {
			collect(children[idx])
		}
	}
	collect(0)

	writeInt16(writer, int16(len(finalOrder)))
	writer.WriteByte('\n')
	for i, id := range finalOrder {
		if i > 0 {
			writer.WriteByte(' ')
		}
		writeInt16(writer, id)
	}
	writer.WriteByte('\n')
}
