package logger

import "math"

type ValueType uint8

const (
	StringType ValueType = iota
	IntType
	Int32Type
	Int64Type
	Float32Type
	Float64Type
	BoolType
	StructType
	ArrayMarshalType
	ArrayType
	AnyType
)

type Record struct {
	Message string
	Level   Level
	KVs     []KV
}

type KV struct {
	Key   string
	Value Value
}

type Value struct {
	String    string
	Int       int64
	Interface interface{}
	ValType   ValueType
}

// ArrayMarshal is implemented by types that can serialize themselves
// directly into a []byte buffer, bypassing encoding/json.
type ArrayMarshal interface {
	MarshalArray(b []byte) ([]byte, error)
}

func AddString(key string, value string) KV {
	return KV{Key: key, Value: Value{String: value, ValType: StringType}}
}

func AddInt(key string, value int) KV {
	return KV{Key: key, Value: Value{Int: int64(value), ValType: IntType}}
}

func AddInt64(key string, value int64) KV {
	return KV{Key: key, Value: Value{Int: value, ValType: Int64Type}}
}

func AddInt32(key string, value int32) KV {
	return KV{Key: key, Value: Value{Int: int64(value), ValType: Int32Type}}
}

func AddFloat32(key string, value float32) KV {
	return KV{Key: key, Value: Value{Int: int64(math.Float32bits(value)), ValType: Float32Type}}
}

func AddFloat64(key string, value float64) KV {
	return KV{Key: key, Value: Value{Int: int64(math.Float64bits(value)), ValType: Float64Type}}
}

func AddBool(key string, value bool) KV {
	v := int64(0)
	if value {
		v = 1
	}
	return KV{Key: key, Value: Value{Int: v, ValType: BoolType}}
}

func AddStruct(key string, value any) KV {
	return KV{Key: key, Value: Value{Interface: value, ValType: StructType}}
}

func AddArray(key string, value any) KV {
	return KV{Key: key, Value: Value{Interface: value, ValType: ArrayType}}
}

func AddArrayMarshal(key string, value ArrayMarshal) KV {
	return KV{Key: key, Value: Value{Interface: value, ValType: ArrayMarshalType}}
}
