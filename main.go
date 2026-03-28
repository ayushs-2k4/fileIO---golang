package main

import (
	"fileIO/logger"
	"fileIO/models"
	"fileIO/writer"
	"fmt"
	"strconv"
	"sync"
	"time"
)

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
			record := logger.Record{
				Message: "Ayush Singhal, " + strconv.Itoa(i),
				Level:   logger.Warn,
				KVs: []logger.KV{
					logger.AddString("my-key", "my-value"),
					logger.AddString("my-key-2", "my-value-2"),
					logger.AddInt64("my-int-key", 34),
					logger.AddStruct("person", models.Person{
						Name: "Ayush Singhal",
						Age:  22,
						Contact: models.ContactInfo{
							Email: "ayush@example.com",
							Phone: "+91-9876543210",
							Social: models.SocialMedia{
								Twitter:  "@ayush",
								LinkedIn: "linkedin.com/in/ayush",
								Stats: models.SocialStats{
									Followers: 4200,
									Posts:     318,
									Verified:  false,
								},
							},
						},
						Address: models.Address{
							Street:  "42 MG Road",
							City:    "Bangalore",
							Country: "India",
							ZipCode: "560001",
							Region: models.Region{
								State:    "Karnataka",
								TimeZone: "Asia/Kolkata",
							},
							Coordinates: models.Coordinates{
								Latitude:  12.9716,
								Longitude: 77.5946,
							},
						},
						Employment: models.Employment{
							Company:    "Blinkit",
							Role:       "Software Engineer",
							Experience: 3,
							Skills:     []string{"Go", "Distributed Systems", "Kafka"},
							Manager: models.Manager{
								Name: "Rahul Sharma",
								Contact: models.ContactInfo{
									Email: "rahul@example.com",
									Phone: "+91-9123456789",
									Social: models.SocialMedia{
										Twitter:  "@rahul",
										LinkedIn: "linkedin.com/in/rahul",
										Stats: models.SocialStats{
											Followers: 12000,
											Posts:     540,
											Verified:  true,
										},
									},
								},
							},
							Salary: models.Salary{
								Total:    2500000,
								Currency: "INR",
								Breakdown: models.SalaryBreakdown{
									Base:  2000000,
									Bonus: 500000,
									TaxRegion: models.TaxRegion{
										Code: "IN-KA",
										Rate: 0.30,
									},
								},
							},
							WorkHistory: models.WorkHistory{
								{Company: "Zomato", Role: "Backend Engineer", YearsExp: 1},
								{Company: "magicpin", Role: "Junior Engineer", YearsExp: 1},
							},
						},
					}),
				},
			}
			jsonEncoder := logger.GetJSONEncoder()
			encodedData, _ := jsonEncoder.Encode(record)
			multiWriter.Write(encodedData)
			logger.PutJSONEncoder(jsonEncoder)
		}()

		wg.Wait()

		fmt.Println()
		fmt.Println(time.Since(st))

		multiWriter.Close()
	}
}
