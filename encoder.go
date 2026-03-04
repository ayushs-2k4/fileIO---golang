package main

type Encoder interface {
	Encode(msg string) ([]byte, error)
}
