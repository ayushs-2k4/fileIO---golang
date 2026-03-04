package main

import (
	"fileIO/writer"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Record struct {
	Message string
	KVs     []KV
}

type KV struct {
	Key   string
	Value interface{}
}

func main() {
	filename := "my-file.txt"
	fileWriter := writer.NewFileWriter(filename)
	consoleWriter := writer.NewConsoleWriter()
	multiWriter := writer.NewMultiWriter(fileWriter, consoleWriter)

	st := time.Now()

	wg := sync.WaitGroup{}

	for i := 0; i < 1; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			record := Record{
				Message: "Ayush Singhal, " + strconv.Itoa(i),
				KVs: []KV{
					{
						Key:   "my-key",
						Value: "my-value",
					},
					{
						Key:   "my-key-2",
						Value: "my-value-2",
					},
				},
			}
			jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
			encodedData, _ := jsonEncoder.Encode(record)
			multiWriter.Write(encodedData)
			_jsonPOOL.Put(jsonEncoder)
		}()
	}

	wg.Wait()

	fmt.Println()
	fmt.Println(time.Since(st))

	multiWriter.Close()
}
