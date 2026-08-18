package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func strUnderSample(str, sample string) string {
	newStr := make([]rune, 0, len(sample))
	strIndex := 0
	for _, value := range sample {
		if (value <= '9' && value >= '0') || value == 'X' {
			newStr = append(newStr, rune(str[strIndex]))
			strIndex++
		} else {
			newStr = append(newStr, value)
		}
	}
	return string(newStr)
}

func equalWithSample(str, sample string) bool {
	if len(str) != len(sample) {
		return false
	}
	for i := 0; i < len(str); i++ {
		if str[i] != sample[i] && sample[i] != 'X' {
			return false
		}
	}
	return true
}

func onlyDigits(str string) string {
	res := make([]rune, 0, 16)
	parts := strings.Split(str, " - ")

	for _, value := range parts[0] {
		if (value <= '9' && value >= '0') || value == 'X' {
			res = append(res, value)
		}

	}
	return string(res)
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)

	n, _ := reader.ReadString('\n')
	n = strings.TrimSpace(n)
	N, _ := strconv.Atoi(n)

	numbers := make([]string, 0, N)
	samples := make([]string, 0, N)
	onlyDigitsSample := make([]string, 0, N)
	for i := 0; i < N; i++ {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		numbers = append(numbers, onlyDigits(line))
	}

	k, _ := reader.ReadString('\n')
	k = strings.TrimSpace(k)

	K, _ := strconv.Atoi(k)

	for i := 0; i < K; i++ {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		samples = append(samples, line)
		onlyDigitsSample = append(onlyDigitsSample, onlyDigits(line))
	}

	for _, num := range numbers {
		for idx, sample := range onlyDigitsSample {
			if equalWithSample(num, sample) {
				res := strUnderSample(num, samples[idx])
				fmt.Println(res)
				break
			}
		}
	}
}
