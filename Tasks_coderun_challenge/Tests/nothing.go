package main

import (
	"fmt"
	"math"
)

func main() {
	for rank := 5; rank < 3000; rank++ {
		for lang := 50; lang < 10000; lang++ {
			result := (1 - ((float64(rank) - 1) / float64(lang)))
			result = math.Pow(result, 2.5)
			result = result * 200
			if result > 61.525 && result < 61.535 {
				fmt.Println("Rank:", rank, "Lang:", lang, "Result:", result)
			}

		}
	}
}
