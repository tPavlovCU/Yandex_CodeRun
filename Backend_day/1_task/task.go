package main

import (
	"bufio"
	"fmt"
	"os"
)

func canPlace(wantPlace byte, badSymbol byte, symbols map[byte]int, badSymbols map[byte]int) bool {
	symbols[wantPlace]--
	badSymbols[badSymbol]--
	if (symbols['a']+symbols['b'] >= badSymbols['c']) && (symbols['a']+symbols['c'] >= badSymbols['b']) && (symbols['b']+symbols['c'] >= badSymbols['a']) {
		symbols[wantPlace]++
		badSymbols[badSymbol]++
		return true
	}
	symbols[wantPlace]++
	return false

}
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)

	symbols := make(map[byte]int, 3)
	reader.ReadByte()
	symbols['a'] = 0
	symbols['b'] = 0
	symbols['c'] = 0
	for i := 0; i < int(n); i++ {
		symbol, _ := reader.ReadByte()
		symbols[symbol]++
	}

	result := make([]byte, 0, n)
	reader.ReadByte()

	badSymbols := make(map[byte]int, 3)
	badString := make([]byte, 0, n)
	for i := 0; i < int(n); i++ {
		badSymbol, _ := reader.ReadByte()
		badSymbols[badSymbol]++
		badString = append(badString, badSymbol)
	}

	for _, value := range badString {
		if value == 'a' {
			if canPlace('b', value, symbols, badSymbols) && symbols['b'] > 0 {
				result = append(result, 'b')
				symbols['b']--
				badSymbols['a']--
			} else if canPlace('c', value, symbols, badSymbols) {
				result = append(result, 'c')
				symbols['c']--
				badSymbols['a']--
			} else {
				fmt.Println("NO", string(value))
			}
		}
		if value == 'b' {
			if canPlace('a', value, symbols, badSymbols) && symbols['a'] > 0 {
				result = append(result, 'a')
				symbols['a']--
				badSymbols['b']--
			} else if canPlace('c', value, symbols, badSymbols) {
				result = append(result, 'c')
				symbols['c']--
				badSymbols['b']--
			} else {
				fmt.Println("NO", string(value))
			}
		}
		if value == 'c' {
			if canPlace('a', value, symbols, badSymbols) && symbols['a'] > 0 {
				result = append(result, 'a')
				symbols['a']--
				badSymbols['c']--
			} else if canPlace('b', value, symbols, badSymbols) {
				result = append(result, 'b')
				symbols['b']--
				badSymbols['c']--
			} else {
				fmt.Println("NO", string(value))
			}
		}
	}
	writer.Write(result)
}
