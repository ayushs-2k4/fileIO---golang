package main

import (
	"path"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"
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
	TimeStampKey     = "timestamp"
	CallerKey        = "caller"
)

const (
	shouldPrettify      = true
	shouldAddCallerInfo = true
)

func (j *JSONEncoder) Encode(rec Record) ([]byte, error) {
	j.b = append(j.b, '{')
	j.addNewLine()
	j.currentLevel++
	j.addTabs()
	j.addKeyValue(AddString(MessageKey, rec.Message))

	j.addCharacter(CommaCharacter)
	j.addNewLine()
	j.addTabs()
	j.addKeyValue(AddString(LevelKey, getLevelString(rec.Level)))

	j.addCharacter(CommaCharacter)
	j.addNewLine()
	j.addTabs()
	j.addKey(TimeStampKey)
	j.addCharacter('"')
	j.b = time.Now().UTC().AppendFormat(j.b, time.RFC3339Nano)
	j.addCharacter('"')

	if shouldAddCallerInfo {
		j.addCharacter(CommaCharacter)
		j.addNewLine()
		j.addTabs()
		j.addKey(CallerKey)
		j.addCharacter('"')
		j.addRawCaller()
		j.addCharacter('"')
	}

	for _, kv := range rec.KVs {
		key := kv.Key
		val := kv.Value

		j.addCharacter(CommaCharacter)
		j.addNewLine()
		j.addTabs()

		j.addKeyValue(KV{
			Key:   key,
			Value: val,
		},
		)

	}

	j.addNewLine()
	j.b = append(j.b, '}')
	j.addNewLine()
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

func (j *JSONEncoder) addKeyValue(kv KV) {
	j.addKey(kv.Key)

	switch kv.Value.ValType {
	case StringType:
		j.addString(kv.Value.String)
	case Int64Type:
		j.addInt(kv.Value.Int)
	case StructType:
		j.addStruct(kv.Value.Interface)
	}

}

func (j *JSONEncoder) addKey(key string) {
	j.addString(key)
	j.b = append(j.b, ':')
	if shouldPrettify {
		j.b = append(j.b, ' ')
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
			j.addKeyValue(AddString(fieldTyp.Name, fieldVal.String()))

		case reflect.Int64:
			j.addKeyValue(AddInt(fieldTyp.Name, fieldVal.Int()))

		case reflect.Struct:
			j.addKeyValue(AddStruct(fieldTyp.Name, fieldVal.Interface()))
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

func (j *JSONEncoder) addRawCaller() {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}

	file = path.Base(file)

	fn := runtime.FuncForPC(pc)
	funcName := "unknown"

	if fn != nil {
		funcName = path.Base(fn.Name())
	}

	j.b = append(j.b, file...)
	j.b = append(j.b, ':')
	j.b = strconv.AppendInt(j.b, int64(line), 10)
	j.b = append(j.b, ' ')
	j.b = append(j.b, funcName...)
}
