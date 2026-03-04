package main

import (
	"sync"
)

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

const (
	NewLineCharacter = '\n'
	TabCharacter     = '\t'
	MessageKey       = "message"
)

func (j *JSONEncoder) Encode(rec Record) ([]byte, error) {
	j.b = append(j.b, '{')
	j.addCharacter(NewLineCharacter)
	j.addCharacter(TabCharacter)
	j.addKeyValue(MessageKey, rec.Message)

	for _, kv := range rec.KVs {
		key := kv.Key
		val := kv.Value

		j.addCharacter(NewLineCharacter)
		j.addCharacter(TabCharacter)

		j.addKeyValue(key, val.(string))
	}

	j.addCharacter(NewLineCharacter)
	j.b = append(j.b, '}')

	res := j.b
	j.b = j.b[:0]

	return res, nil
}

func (j *JSONEncoder) addCharacter(c rune) {
	j.b = append(j.b, byte(c))
}

func (j *JSONEncoder) addKeyValue(key string, value string) {
	j.addString(key)
	j.b = append(j.b, ':')
	j.b = append(j.b, ' ')
	j.addString(value)
}

func (j *JSONEncoder) addString(str string) {
	j.b = append(j.b, '"')
	j.b = append(j.b, []byte(str)...)
	j.b = append(j.b, '"')
}
