package logger

import (
	"fileIO/writer"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ── Benchmark fixture types ───────────────────────────────────────────────────
// Defined here so the logger package has no dependency on package main.

type WorkEntry struct {
	Company  string
	Role     string
	YearsExp int
}

// WorkHistory implements ArrayMarshal for zero-alloc JSON serialization.
type WorkHistory []WorkEntry

func (w WorkHistory) MarshalArray(b []byte) ([]byte, error) {
	b = append(b, '[')
	for i, e := range w {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, `{"company":"`...)
		b = append(b, e.Company...)
		b = append(b, `","role":"`...)
		b = append(b, e.Role...)
		b = append(b, `","years":`...)
		b = strconv.AppendInt(b, int64(e.YearsExp), 10)
		b = append(b, '}')
	}
	b = append(b, ']')
	return b, nil
}

type SocialStats struct {
	Followers int64
	Posts     int64
	Verified  bool
}

type SocialMedia struct {
	Twitter  string
	LinkedIn string
	Stats    SocialStats
}

type ContactInfo struct {
	Email  string
	Phone  string
	Social SocialMedia
}

type Region struct {
	State    string
	TimeZone string
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type Address struct {
	Street      string
	City        string
	Country     string
	ZipCode     string
	Region      Region
	Coordinates Coordinates
}

type TaxRegion struct {
	Code string
	Rate float64
}

type SalaryBreakdown struct {
	Base      float64
	Bonus     float64
	TaxRegion TaxRegion
}

type Salary struct {
	Total     float64
	Currency  string
	Breakdown SalaryBreakdown
}

type Manager struct {
	Name    string
	Contact ContactInfo
}

type Employment struct {
	Company     string
	Role        string
	Experience  int
	Skills      []string
	Manager     Manager
	Salary      Salary
	WorkHistory WorkHistory
}

type Person struct {
	Name       string
	Age        int64
	Contact    ContactInfo
	Address    Address
	Employment Employment
}

// ── Benchmarks ────────────────────────────────────────────────────────────────

func BenchmarkEncoderWriter(b *testing.B) {
	b.ReportAllocs()
	fileWriter := writer.NewFileWriter("bench.log")
	multiWriter := writer.NewMultiWriter(fileWriter)

	record := Record{
		Message: "Ayush Singhal",
		Level:   Warn,
		KVs: []KV{
			AddString("my-key", "my-value"),
			AddString("my-key-2", "my-value-2"),
			AddInt64("my-int-key", 34),
			AddStruct("person", Person{
				Name: "Ayush Singhal",
				Age:  22,
				Contact: ContactInfo{
					Email: "ayush@example.com",
					Phone: "+91-9876543210",
					Social: SocialMedia{
						Twitter: "@ayush",
						Stats:   SocialStats{Followers: 4200, Verified: false},
					},
				},
				Employment: Employment{
					Company: "Blinkit",
					Role:    "Software Engineer",
					Salary: Salary{
						Total:    2500000,
						Currency: "INR",
						Breakdown: SalaryBreakdown{
							Base:      2000000,
							Bonus:     500000,
							TaxRegion: TaxRegion{Code: "IN-KA", Rate: 0.30},
						},
					},
					WorkHistory: WorkHistory{
						{Company: "Zomato", Role: "Junior Engineer", YearsExp: 1},
						{Company: "Swiggy", Role: "Backend Engineer", YearsExp: 2},
					},
				},
			}),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		enc := _jsonPOOL.Get().(*JSONEncoder)
		encodedData, _ := enc.Encode(record)
		multiWriter.Write(encodedData)
		_jsonPOOL.Put(enc)
	}

	b.StopTimer()
	multiWriter.Close()
}

func BenchmarkEncoder(b *testing.B) {
	b.ReportAllocs()
	rec := Record{
		Message: "Ayush Singhal",
		Level:   Warn,
		KVs: []KV{
			AddString("my-key", "my-value"),
			AddInt64("my-int-key", 34),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		enc := _jsonPOOL.Get().(*JSONEncoder)
		enc.Encode(rec) //nolint:errcheck
		_jsonPOOL.Put(enc)
	}
}

func BenchmarkFileWriter(b *testing.B) {
	b.ReportAllocs()
	fileWriter := writer.NewFileWriter("bench.log")

	record := Record{
		Message: "Ayush Singhal",
		Level:   Warn,
		KVs: []KV{
			AddString("my-key", "my-value"),
			AddString("my-key-2", "my-value-2"),
			AddInt64("my-int-key", 34),
			AddStruct("person", Person{
				Name: "Ayush Singhal",
				Age:  22,
				Contact: ContactInfo{
					Email: "ayush@example.com",
					Social: SocialMedia{
						Twitter: "@ayush",
						Stats:   SocialStats{Followers: 4200, Verified: false},
					},
				},
				Address: Address{
					City:    "Bangalore",
					Country: "India",
					Region:  Region{State: "Karnataka", TimeZone: "Asia/Kolkata"},
				},
			}),
		},
	}

	enc := _jsonPOOL.Get().(*JSONEncoder)
	encodedData, _ := enc.Encode(record)
	_jsonPOOL.Put(enc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fileWriter.Write(encodedData)
	}

	b.StopTimer()
	fileWriter.Close()
}

// TestJSONEncoderMethodAllocs measures allocs/op for each individual method of JSONEncoder.
// Run with: go test -v -run TestJSONEncoderMethodAllocs ./logger/
func TestJSONEncoderMethodAllocs(t *testing.T) {
	enc := NewJSONEncoder()

	printAllocs := func(name string, n float64) {
		fmt.Printf("%-40s %.0f allocs/op\n", name, n)
	}

	fmt.Println("--- JSONEncoder method-level allocs/op ---")

	printAllocs("addString", testing.AllocsPerRun(100, func() {
		enc.addString("hello-world")
		enc.reset()
	}))

	printAllocs("addInt", testing.AllocsPerRun(100, func() {
		enc.addInt(12345)
		enc.reset()
	}))

	printAllocs("addRawCaller", testing.AllocsPerRun(100, func() {
		enc.addRawCaller()
		enc.reset()
	}))

	printAllocs("time.Now().UTC().AppendFormat(enc.b, time.RFC3339Nano)", testing.AllocsPerRun(100, func() {
		enc.b = time.Now().UTC().AppendFormat(enc.b, time.RFC3339Nano)
		enc.reset()
	}))

	printAllocs("addKeyValue (string)", testing.AllocsPerRun(100, func() {
		enc.addKeyValue(AddString("k", "v"))
		enc.reset()
	}))

	printAllocs("addKeyValue (int64)", testing.AllocsPerRun(100, func() {
		enc.addKeyValue(AddInt64("k", int64(42)))
		enc.reset()
	}))

	type FlatStruct struct {
		Name string
		Age  int64
	}
	printAllocs("addStruct (flat)", testing.AllocsPerRun(100, func() {
		enc.addStruct(FlatStruct{Name: "Ayush", Age: 22})
		enc.reset()
	}))

	printAllocs("addStruct (nested)", testing.AllocsPerRun(100, func() {
		enc.addStruct(Person{
			Name: "Ayush",
			Age:  22,
			Contact: ContactInfo{
				Email: "ayush@example.com",
				Social: SocialMedia{
					Twitter: "@ayush",
					Stats:   SocialStats{Followers: 4200, Verified: false},
				},
			},
		})
		enc.reset()
	}))

	rec := Record{
		Message: "Ayush Singhal",
		Level:   Warn,
		KVs: []KV{
			AddString("my-key", "my-value"),
			AddString("my-key-2", "my-value-2"),
			AddInt64("my-int-key", 34),
			AddStruct("person", Person{
				Name: "Ayush",
				Age:  22,
				Contact: ContactInfo{
					Email: "ayush@example.com",
					Social: SocialMedia{
						Twitter: "@ayush",
						Stats:   SocialStats{Followers: 4200, Verified: false},
					},
				},
				Employment: Employment{
					Company: "Blinkit",
					Role:    "Software Engineer",
					Salary: Salary{
						Total:    2500000,
						Currency: "INR",
						Breakdown: SalaryBreakdown{
							Base:      2000000,
							Bonus:     500000,
							TaxRegion: TaxRegion{Code: "IN-KA", Rate: 0.30},
						},
					},
					WorkHistory: WorkHistory{
						{Company: "Zomato", Role: "Junior Engineer", YearsExp: 1},
						{Company: "Swiggy", Role: "Backend Engineer", YearsExp: 2},
					},
				},
			}),
		},
	}
	fmt.Println()
	printAllocs("Encode (with pool)", testing.AllocsPerRun(100, func() {
		e := _jsonPOOL.Get().(*JSONEncoder)
		e.Encode(rec) //nolint:errcheck
		_jsonPOOL.Put(e)
	}))
	printAllocs("Encode (no pool)", testing.AllocsPerRun(100, func() {
		e := NewJSONEncoder()
		e.Encode(rec) //nolint:errcheck
	}))
}

func BenchmarkMyLogger10Fields(b *testing.B) {
	b.ReportAllocs()
	w := &writer.DiscardWriter{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		record := Record{
			Message: "test message",
			Level:   Info,
			KVs: []KV{
				AddString("field1", "value1"),
				AddString("field2", "value2"),
				AddString("field3", "value3"),
				AddString("field4", "value4"),
				AddString("field5", "value5"),
				AddString("field6", "value6"),
				AddString("field7", "value7"),
				AddString("field8", "value8"),
				AddString("field9", "value9"),
				AddInt64("field10", 42),
			},
		}
		enc := _jsonPOOL.Get().(*JSONEncoder)
		data, _ := enc.Encode(record)
		w.Write(data)
		_jsonPOOL.Put(enc)
	}
}

func BenchmarkMyLogger10FieldsCreatingOnce(b *testing.B) {
	b.ReportAllocs()
	w := &writer.DiscardWriter{}

	record := Record{
		Message: "test message",
		Level:   Info,
		KVs: []KV{
			AddString("field1", "value1"),
			AddString("field2", "value2"),
			AddString("field3", "value3"),
			AddString("field4", "value4"),
			AddString("field5", "value5"),
			AddString("field6", "value6"),
			AddString("field7", "value7"),
			AddString("field8", "value8"),
			AddString("field9", "value9"),
			AddInt64("field10", 42),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		enc := _jsonPOOL.Get().(*JSONEncoder)
		data, _ := enc.Encode(record)
		w.Write(data)
		_jsonPOOL.Put(enc)
	}
}

func BenchmarkZap10Fields(b *testing.B) {
	b.ReportAllocs()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	core := zapcore.NewCore(encoder, zapcore.AddSync(io.Discard), zap.InfoLevel)
	logger := zap.New(core)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("test message",
			zap.String("field1", "value1"),
			zap.String("field2", "value2"),
			zap.String("field3", "value3"),
			zap.String("field4", "value4"),
			zap.String("field5", "value5"),
			zap.String("field6", "value6"),
			zap.String("field7", "value7"),
			zap.String("field8", "value8"),
			zap.String("field9", "value9"),
			zap.Int("field10", 42),
		)
	}
}

func BenchmarkZap10FieldsCreatingOnce(b *testing.B) {
	b.ReportAllocs()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	core := zapcore.NewCore(encoder, zapcore.AddSync(io.Discard), zap.InfoLevel)
	logger := zap.New(core)

	fields := []zap.Field{
		zap.String("field1", "value1"),
		zap.String("field2", "value2"),
		zap.String("field3", "value3"),
		zap.String("field4", "value4"),
		zap.String("field5", "value5"),
		zap.String("field6", "value6"),
		zap.String("field7", "value7"),
		zap.String("field8", "value8"),
		zap.String("field9", "value9"),
		zap.Int("field10", 42),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("test message", fields...)
	}
}
