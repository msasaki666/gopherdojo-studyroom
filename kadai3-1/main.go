package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"
)

func main() {
	dur := time.Minute * 1
	fmt.Println("制限時間", dur)
	timer := time.NewTimer(dur)
	ch := make(chan bool)
	go textMatching(ch)
	var correctCount int
	for {
		select {
		case <-timer.C:
			fmt.Println("正解数")
			fmt.Println(correctCount)
			os.Exit(0)
		case ok := <-ch:
			if ok {
				correctCount++
			}
			go textMatching(ch)
		}
	}
}

func textMatching(ch chan bool) {
	t := generateWord()
	fmt.Println(t)
	fmt.Println("please input")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ch <- scanner.Text() == t
}

func generateWord() string {
	digit := 3
	b := make([]byte, digit)
	rand.Read(b)

	return base64.StdEncoding.EncodeToString(b)
}
