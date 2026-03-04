package main

type Writer interface {
	Write(b []byte)
	Close()
}
