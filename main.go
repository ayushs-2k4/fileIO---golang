package main

import (
	"fmt"
	"time"
)

const channelSize = 100

func main() {
	filename := "my-file.txt"
	fileWriter := NewFileWriter(filename)

	st := time.Now()

	for i := 0; i < 30000; i++ {
		i := i
		go func() {
			data := fmt.Sprintf("\nAyush Singhal, %d", i)
			jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
			encodedData, _ := jsonEncoder.Encode(data)
			fileWriter.Log(encodedData)
		}()
	}

	fmt.Println(time.Since(st))

	fileWriter.Close()
}
