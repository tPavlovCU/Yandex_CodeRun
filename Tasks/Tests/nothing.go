package main

import (
	"fmt"
)

func main() {
	mapa := make(map[int]int)
	mapa[1] += 1
	fmt.Println(mapa)
}
