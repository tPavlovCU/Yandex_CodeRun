package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type DayResponse struct {
	DayID            string       `json:"day_id"`
	TimeLimitMinutes int          `json:"time_limit_minutes"`
	Budget           int          `json:"budget"`
	ScoreWeights     ScoreWeights `json:"score_weights"`
}

type ScoreWeights struct {
	Rating            float64 `json:"rating"`
	Diversity         float64 `json:"diversity"`
	TravelPenalty     float64 `json:"travel_penalty"`
	BadWeatherPenalty float64 `json:"bad_weather_penalty"`
}

type PlacesResponse struct {
	Places []PlaceBrief `json:"places"`
}

type PlaceBrief struct {
	ID int `json:"id"`
}

type BatchRequest struct {
	IDs []int `json:"ids"`
}

type BatchResponse struct {
	Places []PlaceDetails `json:"places"`
}

type PlaceDetails struct {
	ID               int     `json:"id"`
	Type             string  `json:"type"`
	VisitTimeMinutes int     `json:"visit_time_minutes"`
	Cost             int     `json:"cost"`
	Rating           float64 `json:"rating"`
	WeatherSensitive bool    `json:"weather_sensitive"`
	OpenFrom         string  `json:"open_from"`
	OpenTo           string  `json:"open_to"`
}

type MatrixRequest struct {
	IDs []string `json:"ids"`
}

type MatrixResponse struct {
	IDs    []string `json:"ids"`
	Matrix [][]int  `json:"matrix"`
}

type WeatherResponse struct {
	Weather string `json:"weather"`
}

type InternalPlace struct {
	Idx              int
	ID               int
	Type             string
	TypeID           int
	VisitTimeMinutes int
	Cost             int
	Rating           float64
	IsBadWeather     bool
	OpenFromMins     int
	OpenToMins       int
}

var bestScore float64 = -1e18
var bestPath []int

func parseTime(s string) int {
	if len(s) != 5 || s[2] != ':' {
		return 0
	}
	h, _ := strconv.Atoi(s[0:2])
	m, _ := strconv.Atoi(s[3:5])
	return h*60 + m
}

func main() {
	var baseURL string
	var dayIDStr string

	scanner := os.NewFile(0, "stdin")
	var buf [4096]byte
	n, _ := scanner.Read(buf[:])
	input := string(buf[:n])

	lines := []string{}
	curr := ""
	for i := 0; i < len(input); i++ {
		if input[i] == '\n' || input[i] == '\r' {
			if len(curr) > 0 {
				lines = append(lines, curr)
				curr = ""
			}
		} else {
			curr += string(input[i])
		}
	}
	if len(curr) > 0 {
		lines = append(lines, curr)
	}

	if len(lines) < 2 {
		fmt.Println(0)
		return
	}
	baseURL = lines[0]
	dayIDStr = lines[1]

	client := &http.Client{}

	respDay, err := client.Get(baseURL + "/day/" + dayIDStr)
	if err != nil {
		fmt.Println(0)
		return
	}
	defer respDay.Body.Close()
	var dayData DayResponse
	json.NewDecoder(respDay.Body).Decode(&dayData)

	respBrief, err := client.Get(baseURL + "/places?day_id=" + dayIDStr)
	if err != nil {
		fmt.Println(0)
		return
	}
	defer respBrief.Body.Close()
	var briefData PlacesResponse
	json.NewDecoder(respBrief.Body).Decode(&briefData)

	ids := make([]int, len(briefData.Places))
	for i, p := range briefData.Places {
		ids[i] = p.ID
	}

	batchReqBody, _ := json.Marshal(BatchRequest{IDs: ids})
	respBatch, err := client.Post(baseURL+"/places/batch", "application/json", bytes.NewBuffer(batchReqBody))
	if err != nil {
		fmt.Println(0)
		return
	}
	defer respBatch.Body.Close()
	var batchData BatchResponse
	json.NewDecoder(respBatch.Body).Decode(&batchData)

	matrixReqIDs := make([]string, len(ids)+1)
	matrixReqIDs[0] = "start"
	for i, id := range ids {
		matrixReqIDs[i+1] = strconv.Itoa(id)
	}
	matrixReqBody, _ := json.Marshal(MatrixRequest{IDs: matrixReqIDs})
	respMatrix, err := client.Post(baseURL+"/travel-time/matrix", "application/json", bytes.NewBuffer(matrixReqBody))
	if err != nil {
		fmt.Println(0)
		return
	}
	defer respMatrix.Body.Close()
	var matrixData MatrixResponse
	json.NewDecoder(respMatrix.Body).Decode(&matrixData)

	idToMatIdx := make(map[string]int)
	for i, idStr := range matrixData.IDs {
		idToMatIdx[idStr] = i
	}

	typeMap := make(map[string]int)
	typeCount := 0

	places := make([]InternalPlace, len(batchData.Places))
	for i, p := range batchData.Places {
		respW, err := client.Get(baseURL + "/weather?day_id=" + dayIDStr + "&place_id=" + strconv.Itoa(p.ID))
		isBad := false
		if err == nil {
			var wData WeatherResponse
			json.NewDecoder(respW.Body).Decode(&wData)
			respW.Body.Close()
			if p.WeatherSensitive && (wData.Weather == "rainy" || wData.Weather == "windy") {
				isBad = true
			}
		}

		tID, exists := typeMap[p.Type]
		if !exists {
			tID = typeCount
			typeMap[p.Type] = tID
			typeCount++
		}

		places[i] = InternalPlace{
			Idx:              i,
			ID:               p.ID,
			Type:             p.Type,
			TypeID:           tID,
			VisitTimeMinutes: p.VisitTimeMinutes,
			Cost:             p.Cost,
			Rating:           p.Rating,
			IsBadWeather:     isBad,
			OpenFromMins:     parseTime(p.OpenFrom),
			OpenToMins:       parseTime(p.OpenTo),
		}
	}

	N := len(places)
	dist := make([][]int, N+1)
	for i := 0; i <= N; i++ {
		dist[i] = make([]int, N+1)
	}

	startMatIdx := idToMatIdx["start"]
	for j := 0; j < N; j++ {
		pMatIdx := idToMatIdx[strconv.Itoa(places[j].ID)]
		dist[0][j+1] = matrixData.Matrix[startMatIdx][pMatIdx]
		dist[j+1][0] = matrixData.Matrix[pMatIdx][startMatIdx]
	}
	for i := 0; i < N; i++ {
		p1MatIdx := idToMatIdx[strconv.Itoa(places[i].ID)]
		for j := 0; j < N; j++ {
			p2MatIdx := idToMatIdx[strconv.Itoa(places[j].ID)]
			dist[i+1][j+1] = matrixData.Matrix[p1MatIdx][p2MatIdx]
		}
	}

	weights := dayData.ScoreWeights
	budgetLimit := dayData.Budget
	timeLimit := dayData.TimeLimitMinutes

	var currentPath []int
	var typeUsage uint32

	// ИСПРАВЛЕНО: добавлена переменная travelSum для честного расчета штрафа только за дорогу
	var solve func(currIdx, currTime, travelSum, currBudget, badWeatherCount int, currentRatingSum float64, visited uint32, typesCount int)
	solve = func(currIdx, currTime, travelSum, currBudget, badWeatherCount int, currentRatingSum float64, visited uint32, typesCount int) {
		score := currentRatingSum*weights.Rating + float64(typesCount)*weights.Diversity - float64(travelSum)*weights.TravelPenalty - float64(badWeatherCount)*weights.BadWeatherPenalty
		if score > bestScore && len(currentPath) > 0 {
			bestScore = score
			bestPath = append([]int(nil), currentPath...)
		}

		for i := 0; i < N; i++ {
			if (visited & (1 << i)) != 0 {
				continue
			}
			p := &places[i]
			if currBudget+p.Cost > budgetLimit {
				continue
			}

			travel := dist[currIdx][i+1]
			nextTime := currTime + travel
			if nextTime < p.OpenFromMins {
				nextTime = p.OpenFromMins
			}

			if nextTime+p.VisitTimeMinutes > p.OpenToMins || nextTime+p.VisitTimeMinutes > timeLimit {
				continue
			}

			nextBadWeather := badWeatherCount
			if p.IsBadWeather {
				nextBadWeather++
			}

			newType := false
			if (typeUsage & (1 << p.TypeID)) == 0 {
				newType = true
				typeUsage |= (1 << p.TypeID)
				typesCount++
			}

			currentPath = append(currentPath, p.ID)
			// ИСПРАВЛЕНО: передаем travelSum + travel, сохраняя чистый штраф дороги
			solve(i+1, nextTime+p.VisitTimeMinutes, travelSum+travel, currBudget+p.Cost, nextBadWeather, currentRatingSum+p.Rating, visited|(1<<i), typesCount)
			currentPath = currentPath[:len(currentPath)-1]

			if newType {
				typeUsage &= ^(1 << p.TypeID)
				typesCount--
			}
		}
	}

	solve(0, 0, 0, 0, 0, 0, 0, 0)

	if len(bestPath) == 0 {
		fmt.Println(0)
	} else {
		fmt.Println(len(bestPath))
		for i, id := range bestPath {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(id)
		}
		fmt.Println()
	}
}
