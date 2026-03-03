package main

import (
	buffer2 "fileIO/buffer"
	"time"
)

func main() {
	fileName := "my-file.txt"
	buffer, err := buffer2.NewBuffer(fileName)
	if err != nil {
		panic(err)
	}
	defer buffer.Sync()
	for i := 0; i < 1000; i++ {
		go buffer.Write([]byte("\nDevansh Singhal"))
	}
	time.Sleep(100 * time.Millisecond)
}
