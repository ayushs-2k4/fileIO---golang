package main

import (
	"reflect"
	"strconv"
	"sync"
)

var _jsonPOOL = sync.Pool{New: func() any {
	return NewJSONEncoder()
}}

type JSONEncoder struct {
	currentLevel int
	b            []byte
}

func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{
		currentLevel: 0,
		b:            make([]byte, 0, 1024),
	}
}

const (
	NewLineCharacter = '\n'
	TabCharacter     = '\t'
	CommaCharacter   = ','
	MessageKey       = "message"
	LevelKey         = "level"
)

const (
	shouldPrettify = false
)

func (j *JSONEncoder) Encode(rec Record) ([]byte, error) {
	j.b = append(j.b, '{')
	j.addNewLine()
	j.currentLevel++
	j.addTabs()
	j.addKeyValue(MessageKey, Value{
		val:     rec.Message,
		valType: reflect.String,
	})

	j.addCharacter(CommaCharacter)
	j.addNewLine()
	j.addTabs()
	j.addKeyValue(LevelKey, Value{
		val:     getLevelString(rec.Level),
		valType: reflect.String,
	})

	for _, kv := range rec.KVs {
		key := kv.Key
		val := kv.Value

		j.addCharacter(CommaCharacter)
		j.addNewLine()
		j.addTabs()

		j.addKeyValue(key, *val)

	}

	j.addNewLine()
	j.b = append(j.b, '}')
	j.currentLevel--

	res := j.b
	j.reset()

	return res, nil
}

func (j *JSONEncoder) addNewLine() {
	if shouldPrettify {
		j.addCharacter(NewLineCharacter)
	}
}

func (j *JSONEncoder) addTabs() {
	if shouldPrettify {
		for i := 0; i < j.currentLevel; i++ {
			j.addCharacter(TabCharacter)
		}
	}
}

func (j *JSONEncoder) addCharacter(c rune) {
	j.b = append(j.b, byte(c))
}

func (j *JSONEncoder) addKeyValue(key string, value Value) {
	j.addString(key)
	j.b = append(j.b, ':')
	if shouldPrettify {
		j.b = append(j.b, ' ')
	}

	switch value.valType {
	case reflect.String:
		j.addString(value.val.(string))
	case reflect.Int64:
		j.addInt(value.val.(int64))
	case reflect.Struct:
		j.addStruct(value.val)
	}

}

func getLevelString(level Level) string {
	switch level {
	case Error:
		return "ERROR"
	case Warn:
		return "WARN"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	}

	return "N/A"
}

func (j *JSONEncoder) addString(str string) {
	j.b = append(j.b, '"')
	j.b = append(j.b, str...)
	j.b = append(j.b, '"')
}

func (j *JSONEncoder) addInt(val int64) {
	j.b = strconv.AppendInt(j.b, val, 10)
}

func (j *JSONEncoder) addStruct(value any) {
	val := reflect.ValueOf(value)
	typ := reflect.TypeOf(value)

	j.addCharacter('{')
	j.currentLevel++

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)

		if !fieldTyp.IsExported() {
			continue
		}

		j.addNewLine()
		j.addTabs()

		switch fieldVal.Kind() {
		case reflect.String:
			j.addKeyValue(fieldTyp.Name, Value{
				val:     fieldVal.String(),
				valType: reflect.String,
			})

		case reflect.Int64:
			j.addKeyValue(fieldTyp.Name, Value{
				val:     fieldVal.Int(),
				valType: reflect.Int64,
			})

		case reflect.Struct:
			j.addKeyValue(fieldTyp.Name, Value{
				val:     fieldVal.Interface(),
				valType: reflect.Struct,
			})
		}

		if i < val.NumField()-1 {
			j.addCharacter(CommaCharacter)
		}
	}

	j.addNewLine()
	j.currentLevel--
	j.addTabs()
	j.addCharacter('}')
}

func (j *JSONEncoder) reset() {
	j.b = j.b[:0]
	j.currentLevel = 0
}
