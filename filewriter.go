package main

import (
	"bufio"
	"os"
)

const channelSize = 100

type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
	ch     chan []byte
	done   chan struct{}
}

func NewFileWriter(filename string) *FileWriter {
	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	writer := bufio.NewWriter(file)

	ch := make(chan []byte, channelSize)

	fileWriter := &FileWriter{
		file:   file,
		writer: writer,
		ch:     ch,
		done:   make(chan struct{}),
	}

	go fileWriter.run()

	return fileWriter
}

func (f *FileWriter) run() {
	for msg := range f.ch {
		f.writer.Write(msg)
		//time.Sleep(100 * time.Millisecond)
	}

	// channel closed → flush remaining data
	f.writer.Flush()
	f.file.Close()

	close(f.done)
}

func (f *FileWriter) Write(b []byte) {
	f.ch <- b
}

func (f *FileWriter) Close() {
	close(f.ch) // signal no more logs
	<-f.done
}
