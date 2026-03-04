package main

import "sync"

var _jsonPOOL = sync.Pool{New: func() any {
	return NewJSONEncoder()
}}

type JSONEncoder struct {
	b []byte
}

func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{
		b: make([]byte, 0, 1024),
	}
}

func (j *JSONEncoder) Encode(msg string) ([]byte, error) {
	j.b = append(j.b, '{')
	j.b = append(j.b, []byte(msg)...)
	j.b = append(j.b, '}')

	res := j.b
	j.b = j.b[:0]

	_jsonPOOL.Put(j)

	return res, nil
}
