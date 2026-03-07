package main

import (
	"fileIO/writer"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
)

type Level int

const (
	Debug Level = iota
	Error
	Warn
	Info
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

type ValueType uint8

const (
	StringType = iota
	IntType
	Int32Type
	Int64Type
	Float32Type
	Float64Type
	StructType
)

func AddString(key string, value string) KV {
	return KV{
		Key: key,
		Value: Value{
			String:  value,
			ValType: StringType,
		},
	}
}

func AddInt(key string, value int) KV {
	return KV{
		Key: key,
		Value: Value{
			Int:     int64(value),
			ValType: IntType,
		},
	}
}

func AddInt64(key string, value int64) KV {
	return KV{
		Key: key,
		Value: Value{
			Int:     value,
			ValType: Int64Type,
		},
	}
}

func AddInt32(key string, value int32) KV {
	return KV{
		Key: key,
		Value: Value{
			Int:     int64(value),
			ValType: Int32Type,
		},
	}
}

func AddFloat32(key string, value float32) KV {
	return KV{
		Key: key,
		Value: Value{
			Int:     int64(math.Float32bits(value)),
			ValType: Float32Type,
		},
	}
}

func AddFloat64(key string, value float64) KV {
	return KV{
		Key: key,
		Value: Value{
			Int:     int64(math.Float64bits(value)),
			ValType: Float64Type,
		},
	}
}

func AddStruct(key string, value any) KV {
	return KV{
		Key: key,
		Value: Value{
			Interface: value,
			ValType:   StructType,
		},
	}
}

func main() {
	filename := "my-file.txt"
	fileWriter := writer.NewFileWriter(filename)
	consoleWriter := writer.NewConsoleWriter()
	multiWriter := writer.NewMultiWriter(fileWriter, consoleWriter)

	st := time.Now()

	wg := sync.WaitGroup{}

	for i := 0; i < 1; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			record := Record{
				Message: "Ayush Singhal, " + strconv.Itoa(i),
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
								Twitter:  "@ayush",
								LinkedIn: "linkedin.com/in/ayush",
								Stats: SocialStats{ // level 4
									Followers: 4200,
									Posts:     318,
									Verified:  false,
								},
							},
						},
						Address: Address{
							Street:  "42 MG Road",
							City:    "Bangalore",
							Country: "India",
							ZipCode: "560001",
							Region: Region{ // level 3
								State:    "Karnataka",
								TimeZone: "Asia/Kolkata",
							},
							Coordinates: Coordinates{ // level 3
								Latitude:  12.9716,
								Longitude: 77.5946,
							},
						},
						Employment: Employment{
							Company:    "Blinkit",
							Role:       "Software Engineer",
							Experience: 3,
							Skills:     []string{"Go", "Distributed Systems", "Kafka"},
							Manager: Manager{ // level 3
								Name: "Rahul Sharma",
								Contact: ContactInfo{ // level 4
									Email: "rahul@example.com",
									Phone: "+91-9123456789",
									Social: SocialMedia{
										Twitter:  "@rahul",
										LinkedIn: "linkedin.com/in/rahul",
										Stats: SocialStats{ // level 5
											Followers: 12000,
											Posts:     540,
											Verified:  true,
										},
									},
								},
							},
							Salary: Salary{ // level 3
								Total:    2500000,
								Currency: "INR",
								Breakdown: SalaryBreakdown{ // level 4
									Base:  2000000,
									Bonus: 500000,
									TaxRegion: TaxRegion{ // level 5
										Code: "IN-KA",
										Rate: 0.30,
									},
								},
							},
						},
					}),
				},
			}
			jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
			encodedData, _ := jsonEncoder.Encode(record)
			multiWriter.Write(encodedData)
			_jsonPOOL.Put(jsonEncoder)
		}()
	}

	wg.Wait()

	fmt.Println()
	fmt.Println(time.Since(st))

	multiWriter.Close()
}

type Person struct {
	Name       string
	Age        int64
	Contact    ContactInfo
	Address    Address
	Employment Employment
}

// Level 2
type ContactInfo struct {
	Email  string
	Phone  string
	Social SocialMedia
}

type Address struct {
	Street      string
	City        string
	Country     string
	ZipCode     string
	Region      Region
	Coordinates Coordinates
}

type Employment struct {
	Company    string
	Role       string
	Experience int
	Skills     []string
	Manager    Manager
	Salary     Salary
}

// Level 3
type SocialMedia struct {
	Twitter  string
	LinkedIn string
	Stats    SocialStats
}

type Region struct {
	State    string
	TimeZone string
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type Manager struct {
	Name    string
	Contact ContactInfo // reuses ContactInfo — level 3 → 4 via nesting
}

type Salary struct {
	Total     float64
	Currency  string
	Breakdown SalaryBreakdown
}

// Level 4
type SocialStats struct {
	Followers int64
	Posts     int64
	Verified  bool
}

type SalaryBreakdown struct {
	Base      float64
	Bonus     float64
	TaxRegion TaxRegion
}

// Level 5
type TaxRegion struct {
	Code string
	Rate float64
}
