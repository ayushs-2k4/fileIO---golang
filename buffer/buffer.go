package buffer

import (
	"fmt"
	"io"
	"os"
)

const MaxSize = 8

type Buffer struct {
	writer io.Writer
	data   []byte
}

func NewBuffer(filename string) (*Buffer, error) {
	file, err := createFileIfNotExists(filename)
	if err != nil {
		return nil, err
	}

	return &Buffer{
		writer: file,
		data:   make([]byte, 0, MaxSize),
	}, nil
}

func (b *Buffer) Write(data []byte) (bool, error) {
	b.data = append(b.data, data...)

	if len(b.data) >= MaxSize {
		b.Sync()
	}

	return true, nil
}

func (b *Buffer) Sync() error {
	fmt.Println("Syncing...")
	if len(b.data) == 0 {
		return nil
	}

	n, err := b.writer.Write(b.data)
	if err != nil {
		return err
	}
	b.data = b.data[:0]
	fmt.Println(fmt.Sprintf("Written %d bytes", n))
	return nil
}

func (b *Buffer) Close() error {
	b.Sync()
	if file, ok := b.writer.(*os.File); ok {
		return file.Close()
	}
	return nil
}
