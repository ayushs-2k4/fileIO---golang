package main

import (
	"fileIO/writer"
	"fmt"
	"io"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkEncoderWriter(b *testing.B) {
	b.ReportAllocs()
	fileWriter := writer.NewFileWriter("bench.log")
	//consoleWriter := writer.NewConsoleWriter()
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

	record.Message = "Ayush Singhal"
	for i := 0; i < b.N; i++ {
		jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
		encodedData, _ := jsonEncoder.Encode(record)

		multiWriter.Write(encodedData)

		_jsonPOOL.Put(jsonEncoder)
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
		jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
		jsonEncoder.Encode(rec) //nolint:errcheck
		_jsonPOOL.Put(jsonEncoder)
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

	b.ResetTimer()

	record.Message = "Ayush Singhal"
	jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
	encodedData, _ := jsonEncoder.Encode(record)
	for i := 0; i < b.N; i++ {

		fileWriter.Write(encodedData)

		_jsonPOOL.Put(jsonEncoder)
	}

	b.StopTimer()
	fileWriter.Close()
}

// TestJSONEncoderMethodAllocs measures allocs/op for each individual method of JSONEncoder.
// Run with: go test -v -run TestJSONEncoderMethodAllocs
func TestJSONEncoderMethodAllocs(t *testing.T) {
	enc := NewJSONEncoder()

	printAllocs := func(name string, n float64) {
		fmt.Printf("%-40s %.0f allocs/op\n", name, n)
	}

	fmt.Println("--- JSONEncoder method-level allocs/op ---")

	// addString
	printAllocs("addString", testing.AllocsPerRun(100, func() {
		enc.addString("hello-world")
		enc.reset()
	}))

	// addInt
	printAllocs("addInt", testing.AllocsPerRun(100, func() {
		enc.addInt(12345)
		enc.reset()
	}))

	// addCaller  (runtime.Caller + string concat)
	printAllocs("addRawCaller", testing.AllocsPerRun(100, func() {
		enc.addRawCaller()
		enc.reset()
	}))

	// time.Now().UTC().Format  (isolated — the timestamp line in Encode)
	printAllocs("time.Now().UTC().AppendFormat(enc.b, time.RFC3339Nano)", testing.AllocsPerRun(100, func() {
		enc.b = time.Now().UTC().AppendFormat(enc.b, time.RFC3339Nano)
		enc.reset()
	}))

	// addKeyValue with a string Value
	printAllocs("addKeyValue (string)", testing.AllocsPerRun(100, func() {
		enc.addKeyValue(AddString("k", "v"))
		enc.reset()
	}))

	// addKeyValue with an int64 Value
	printAllocs("addKeyValue (int64)", testing.AllocsPerRun(100, func() {
		enc.addKeyValue(AddInt64("k", int64(42)))
		enc.reset()
	}))

	// addStruct (flat struct — no nested struct)
	type FlatStruct struct {
		Name string
		Age  int64
	}
	printAllocs("addStruct (flat)", testing.AllocsPerRun(100, func() {
		enc.addStruct(FlatStruct{Name: "Ayush", Age: 22})
		enc.reset()
	}))

	// addStruct (nested struct)
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

	// Full Encode — with pool
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
	writer := &writer.DiscardWriter{}

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
		jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)

		data, _ := jsonEncoder.Encode(record)

		writer.Write(data)

		_jsonPOOL.Put(jsonEncoder)
	}
}

func BenchmarkMyLogger10FieldsCreatingOnce(b *testing.B) {
	b.ReportAllocs()
	writer := &writer.DiscardWriter{}

	b.ResetTimer()

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

	for i := 0; i < b.N; i++ {
		jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)

		data, _ := jsonEncoder.Encode(record)

		writer.Write(data)

		_jsonPOOL.Put(jsonEncoder)
	}
}

func BenchmarkZap10Fields(b *testing.B) {
	b.ReportAllocs()
	encoderCfg := zap.NewProductionEncoderConfig()

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(io.Discard),
		zap.InfoLevel,
	)

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

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(io.Discard),
		zap.InfoLevel,
	)

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
