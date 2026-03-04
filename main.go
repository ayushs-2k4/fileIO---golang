package main

import (
	"fmt"
	"time"
)

const channelSize = 100

func main() {
	filename := "my-file.txt"
	fileLogger := NewFileLogger(filename)

	st := time.Now()

	for i := 0; i < 30000; i++ {
		i := i
		go func() {
			data := fmt.Sprintf("\nDevansh Singhal, %d", i)
			//jsonEncoder := NewJSONEncoder()
			jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
			encodedData, _ := jsonEncoder.Encode(data)
			fileLogger.Log(encodedData)
		}()
	}

	fmt.Println(time.Since(st))

	time.Sleep(5 * time.Second)

	for i := 0; i < 300; i++ {
		data := fmt.Sprintf("Devansh Singhal, %d", i)
		jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
		encodedData, _ := jsonEncoder.Encode(data)
		fmt.Println(fmt.Sprintf("producer 2: time: %s, i: %d", time.Now(), i))
		fileLogger.Log(encodedData)
	}

	fileLogger.Close()
}
