package logger

type Encoder interface {
	Encode(rec Record) ([]byte, error)
}
