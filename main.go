package main

import (
	"fmt"
	"os"
)

func main() {
	fileName := "my-file.txt"
	writeToFile(fileName, []byte("Devansh Singhal"))
}

func createFileIfNotExists(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if file == nil || err != nil {
		return nil, err
	}

	return file, nil
}

func writeToFile(filename string, b []byte) (bool, error) {
	file, err := createFileIfNotExists(filename)
	if err != nil || file == nil {
		return false, err
	}
	defer file.Close()

	n, err := file.Write(b)
	if err != nil {
		return false, err
	}
	fmt.Println(fmt.Sprintf("Written %d bytes", n))
	return true, nil
}
