// This is a test example and you can run it to store data constantly.
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tempdb "github.com/tempdb-labs/tempdb-go/lib"
)

func main() {

	client, err := tempdb.NewClient(tempdb.Config{
		Addr: "tempdb1.tempdb.xyz:8081",
		URL:  "<url here>",
	})
	if err != nil {
		log.Fatalf("failed to connect to client: %v", err)
	}
	defer client.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatalf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		log.Fatalf("error reading CSV header: %v", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("error reading CSV records: %v", err)
	}

	for range ticker.C {
		record := records[time.Now().UnixNano()%int64(len(records))]

		product := make(map[string]any)
		for i, value := range record {
			value = strings.TrimSpace(value)
			if value == "" {
				continue
			}

			switch header[i] {
			case "samCID":
				product["customer_id"] = value

			case "Gender":
				product["gender"] = value

			case "Age Group":
				product["age_group"] = value

			case "Purchase Date":
				product["purchase_date"] = value

			case "Product Category":
				product["category"] = value

			case "Discount Availed":
				product["discount_availed"] = value

			case "Discount Name":
				product["discount_name"] = value

			case "Discount Amount (INR)":
				product["discount"] = value

			case "Gross Amount":
				product["price"] = value

			case "Net Amount":
				product["net_amount"] = value

			case "Purchase Method":
				product["payment_method"] = value

			case "Location":
				product["location"] = value
			}
		}

		timestamp := time.Now().UnixNano()
		key := fmt.Sprintf("product_%d", timestamp)

		response, err := client.Store(key, product)
		if err != nil {
			log.Printf("error storing product info: %v", err)
		}
		log.Printf("stored product with response: %v\n", response)
	}
}
