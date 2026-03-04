package main

import (
	"fmt"
	"time"
)

const channelSize = 100

func main() {
	filename := "my-file.txt"
	fileLogger := NewFileLogger(filename)

	for i := 0; i < 300; i++ {
		data := fmt.Sprintf("\nDevansh Singhal, %d", i)
		byteData := []byte(data)
		fmt.Println(fmt.Sprintf("producer 1: time: %s, i: %d", time.Now(), i))
		fileLogger.Log(byteData)
	}

	time.Sleep(5 * time.Second)

	for i := 0; i < 300; i++ {
		data := fmt.Sprintf("\nDevansh Singhal, %d", i)
		byteData := []byte(data)
		fmt.Println(fmt.Sprintf("producer 2: time: %s, i: %d", time.Now(), i))
		fileLogger.Log(byteData)
	}

	fileLogger.Close()
}
