package main

import (
	"fileIO/writer"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type Level int

const (
	Debug Level = iota
	Error
	Warn
	Info
)

type Record struct {
	Message string
	Level   Level
	KVs     []KV
}

type KV struct {
	Key   string
	Value *Value
}

type Value struct {
	val     interface{}
	valType reflect.Kind
}

func AddString(key string, value string) KV {
	return KV{
		Key: key,
		Value: &Value{
			val:     value,
			valType: reflect.String,
		},
	}
}

func AddInt(key string, value int64) KV {
	return KV{
		Key: key,
		Value: &Value{
			val:     value,
			valType: reflect.Int64,
		},
	}
}

func AddStruct(key string, value any) KV {
	return KV{
		Key: key,
		Value: &Value{
			val:     value,
			valType: reflect.Struct,
		},
	}
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
				Level:   Warn,
				KVs: []KV{
					AddString("my-key", "my-value"),
					AddString("my-key-2", "my-value-2"),
					AddInt("my-int-key", 34),
					AddStruct("my-struct-key", MyStruct{
						Name:   "Ayush",
						Age:    22,
						MyInfo: MyInfo{Gender: "Male"},
					}),
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

type MyStruct struct {
	Name   string
	Age    int64
	MyInfo MyInfo
}

type MyInfo struct {
	Gender string
}
