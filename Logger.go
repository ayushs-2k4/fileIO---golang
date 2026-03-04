package main

type Logger interface {
	Log(b []byte)
	Close()
}
