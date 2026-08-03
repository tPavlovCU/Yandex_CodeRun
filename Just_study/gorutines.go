package main

import (
	"fmt"
)

func playerOne(ch chan string) {
	ch <- "Ping"

}

func playerTwo(ch chan string) {
	ping := <-ch
	fmt.Println(ping)

	ch <- "Pong"
}

func main() {
	ch := make(chan string)
	go playerTwo(ch)
	go playerOne(ch)
	pong := <-ch

	fmt.Println(pong)
}
