package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Name     string    `json:"name"`
	Type     string    `json:"Type"`
	Id       int64     `json:"id"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Id           int64  `json:"id"`
	Type         string `json:"type"`
	Date         string `json:"date"`
	DateUnixtime string `json:"date_unixtime"`
	From         string `json:"from"`
	FromId       string `json:"from_id"`
	Text         any    `json:"text"`
}

type Word struct {
	txt     string
	from    string
	printed bool
}

func lower(str string) string {
	res := make([]rune, 0, len(str))
	for _, char := range str {
		if char <= 'Я' && char >= 'А' {
			char += 'а' - 'А'
			res = append(res, char)
		} else if char <= 'Z' && char >= 'A' {
			char += 'a' - 'A'
			res = append(res, char)
		} else if ('A' <= char && char <= 'Z') || ('А' <= char && char <= 'Я') || ('a' <= char && char <= 'z') || ('а' <= char && char <= 'я') {
			res = append(res, char)
		}
	}
	return string(res)
}

func toWords(russianWords map[string]struct{}, RES map[string]Word, line string, from string, file *os.File) {
	parts := strings.Split(line, " ")
	for _, word := range parts {
		lowerW := lower(word)

		w, was := RES[lowerW]
		if !was {
			newWord := Word{lowerW, from, false}
			RES[lowerW] = newWord
		} else {
			if w.from == "Солнышко" && from == "Павлов Тимофей" && w.printed == false {
				_, inRussianWords := russianWords[lowerW]
				if !inRussianWords {
					fmt.Fprintf(file, "%v\n", lowerW)
					w := RES[lowerW]
					w.printed = true
					RES[lowerW] = w
				}
			}
		}
	}

}

func main() {
	file, _ := os.ReadFile("result.json")
	var Data Result
	err := json.Unmarshal(file, &Data)
	if err != nil {
		panic(err)
	}
	messages := Data.Messages

	result, _ := os.Create("res.txt")
	defer result.Close()

	words, _ := os.Open("russian.txt")
	defer words.Close()

	RES := make(map[string]Word, 1000000)
	russianWords := make(map[string]struct{}, 1500000)

	scanner := bufio.NewScanner(words)
	for scanner.Scan() {
		russianWords[scanner.Text()] = struct{}{}
	}

	for _, value := range messages {
		txt := value.Text
		from := value.From
		v, ok := txt.(string)
		if !ok {
			arr := value.Text.([]any)
			str, ok := arr[0].(string)
			if ok {
				toWords(russianWords, RES, str, from, result)
			}
		} else {
			toWords(russianWords, RES, v, from, result)
		}

	}

}
