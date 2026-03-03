package buffer

import (
	"fmt"
	"os"
)

const MaxSize = 80

type Buffer struct {
	file *os.File
	data []byte
}

func NewBuffer(filename string) (*Buffer, error) {
	file, err := createFileIfNotExists(filename)
	if err != nil {
		return nil, err
	}

	return &Buffer{
		file: file,
		data: nil,
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

	n, err := b.file.Write(b.data)
	if err != nil {
		return err
	}
	b.data = b.data[:0]
	fmt.Println(fmt.Sprintf("Written %d bytes", n))
	return nil
}

func (b *Buffer) Close() error {
	b.Sync()
	return b.file.Close()
}
