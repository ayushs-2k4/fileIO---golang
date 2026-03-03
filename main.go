package main

import (
	buffer2 "fileIO/buffer"
	"fmt"
	"sync"
)

func main() {
	fileName := "my-file.txt"
	var wg sync.WaitGroup
	buffer, err := buffer2.NewBuffer(fileName)
	if err != nil {
		panic(err)
	}
	defer buffer.Sync()

	data := "\nDevansh Singhal"
	byteData := []byte(data)
	fmt.Println("Size of data: ", len(byteData))

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			buffer.Write(byteData)
			wg.Done()
		}()
	}
	wg.Wait()
}
