package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	filename := "my-file.txt"
	fileWriter := NewFileWriter(filename)

	st := time.Now()

	wg := sync.WaitGroup{}

	for i := 0; i < 300000; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := "\nAyush Singhal, " + strconv.Itoa(i)
			jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
			encodedData, _ := jsonEncoder.Encode(data)
			fileWriter.Write(encodedData)
			_jsonPOOL.Put(jsonEncoder)
		}()
	}

	wg.Wait()

	fmt.Println(time.Since(st))

	fileWriter.Close()
}
