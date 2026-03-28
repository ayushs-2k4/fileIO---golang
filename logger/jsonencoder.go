package logger

import (
	"encoding/json"
	"math"
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

func GetJSONEncoder() *JSONEncoder {
	return _jsonPOOL.Get().(*JSONEncoder)
}

func PutJSONEncoder(enc *JSONEncoder) {
	_jsonPOOL.Put(enc)
}

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
	newLineCharacter = '\n'
	tabCharacter     = '\t'
	commaCharacter   = ','
	messageKey       = "message"
	levelKey         = "level"
	timeStampKey     = "timestamp"
	callerKey        = "caller"
)

const (
	shouldPrettify      = true
	shouldAddCallerInfo = false
)

func (j *JSONEncoder) Encode(rec Record) ([]byte, error) {
	j.b = append(j.b, '{')
	j.addNewLine()
	j.currentLevel++
	j.addTabs()
	j.addKeyValue(AddString(messageKey, rec.Message))

	j.addLevel(rec)
	j.addTimestamp()

	if shouldAddCallerInfo {
		j.addCallerInfo()
	}

	for _, kv := range rec.KVs {
		j.addCharacter(commaCharacter)
		j.addNewLine()
		j.addTabs()
		j.addKeyValue(kv)
	}

	j.addNewLine()
	j.b = append(j.b, '}')
	j.addNewLine()
	j.currentLevel--

	res := j.b
	j.reset()

	return res, nil
}

func (j *JSONEncoder) addLevel(rec Record) {
	j.addCharacter(commaCharacter)
	j.addNewLine()
	j.addTabs()
	j.addKeyValue(AddString(levelKey, rec.Level.String()))
}

func (j *JSONEncoder) addTimestamp() {
	j.addCharacter(commaCharacter)
	j.addNewLine()
	j.addTabs()
	j.addKey(timeStampKey)
	j.addCharacter('"')
	j.b = time.Now().UTC().AppendFormat(j.b, time.RFC3339Nano)
	j.addCharacter('"')
}

func (j *JSONEncoder) addCallerInfo() {
	j.addCharacter(commaCharacter)
	j.addNewLine()
	j.addTabs()
	j.addKey(callerKey)
	j.addCharacter('"')
	j.addRawCaller()
	j.addCharacter('"')
}

func (j *JSONEncoder) addNewLine() {
	if shouldPrettify {
		j.addCharacter(newLineCharacter)
	}
}

func (j *JSONEncoder) addTabs() {
	if shouldPrettify {
		for i := 0; i < j.currentLevel; i++ {
			j.addCharacter(tabCharacter)
		}
	}
}

func (j *JSONEncoder) addCharacter(c rune) {
	j.b = append(j.b, byte(c))
}

func (j *JSONEncoder) addKeyValue(kv KV) {
	j.addKey(kv.Key)
	j.addValue(kv.Value)
}

func (j *JSONEncoder) addValue(v Value) {
	switch v.ValType {
	case StringType:
		j.addString(v.String)
	case IntType, Int32Type, Int64Type:
		j.addInt(v.Int)
	case Float32Type:
		j.addFloat(float64(math.Float32frombits(uint32(v.Int))))
	case Float64Type:
		j.addFloat(math.Float64frombits(uint64(v.Int)))
	case BoolType:
		j.addBool(v.Int == 1)
	case StructType:
		j.addStruct(v.Interface)
	case ArrayMarshalType:
		j.addArrayMarshal(v.Interface.(ArrayMarshal))
	case ArrayType:
		j.addArray(reflect.ValueOf(v.Interface))
	}
}

func (j *JSONEncoder) addKey(key string) {
	j.addString(key)
	j.b = append(j.b, ':')
	if shouldPrettify {
		j.b = append(j.b, ' ')
	}
}

func (j *JSONEncoder) addAny(val any) {
	b, err := json.Marshal(val)
	if err != nil {
		panic("jsonencoder: failed to marshal value: " + err.Error())
	}
	j.b = append(j.b, b...)
}

func (j *JSONEncoder) addArrayMarshal(arr ArrayMarshal) {
	var err error
	j.b, err = arr.MarshalArray(j.b)
	if err != nil {
		panic("jsonencoder: failed to marshal array: " + err.Error())
	}
}

func (j *JSONEncoder) addArray(arr reflect.Value) {
	b, err := json.Marshal(arr.Interface())
	if err != nil {
		panic("jsonencoder: failed to marshal array: " + err.Error())
	}
	j.b = append(j.b, b...)
}

func (j *JSONEncoder) addString(str string) {
	j.b = append(j.b, '"')
	j.b = append(j.b, str...)
	j.b = append(j.b, '"')
}

func (j *JSONEncoder) addInt(val int64) {
	j.b = strconv.AppendInt(j.b, val, 10)
}

func (j *JSONEncoder) addFloat(val float64) {
	j.b = strconv.AppendFloat(j.b, val, 'f', -1, 64)
}

func (j *JSONEncoder) addBool(val bool) {
	if val {
		j.b = append(j.b, "true"...)
	} else {
		j.b = append(j.b, "false"...)
	}
}

func (j *JSONEncoder) addStruct(value any) {
	val := reflect.ValueOf(value)
	typ := reflect.TypeOf(value)

	j.addCharacter('{')
	j.currentLevel++

	// Use comma-before pattern so trailing commas never appear,
	// even if the last struct field(s) are unexported.
	first := true
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)

		if !fieldTyp.IsExported() {
			continue
		}

		if !first {
			j.addCharacter(commaCharacter)
		}
		first = false

		j.addNewLine()
		j.addTabs()
		j.addReflectionValue(fieldVal, fieldTyp)
	}

	j.addNewLine()
	j.currentLevel--
	j.addTabs()
	j.addCharacter('}')
}

func (j *JSONEncoder) addReflectionValue(fieldVal reflect.Value, fieldTyp reflect.StructField) {
	switch fieldVal.Kind() {
	case reflect.String:
		j.addKeyValue(AddString(fieldTyp.Name, fieldVal.String()))
	case reflect.Int, reflect.Int32, reflect.Int64:
		j.addKeyValue(AddInt64(fieldTyp.Name, fieldVal.Int()))
	case reflect.Float32, reflect.Float64:
		j.addKeyValue(AddFloat64(fieldTyp.Name, fieldVal.Float()))
	case reflect.Bool:
		j.addKeyValue(AddBool(fieldTyp.Name, fieldVal.Bool()))
	case reflect.Struct:
		j.addKeyValue(AddStruct(fieldTyp.Name, fieldVal.Interface()))
	case reflect.Array, reflect.Slice:
		if am, ok := fieldVal.Interface().(ArrayMarshal); ok {
			j.addKeyValue(AddArrayMarshal(fieldTyp.Name, am))
		} else {
			j.addKeyValue(AddArray(fieldTyp.Name, fieldVal.Interface()))
		}
	}
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
