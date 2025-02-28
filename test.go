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
	// Connect1()

	client, err := tempdb.NewClient(tempdb.Config{
		Addr: "0.0.0.0:8081",
		URL:  "tempdb://admin:123456789@workspace:cb4552273c5c/ecommerce",
	})
	if err != nil {
		log.Fatalf("failed to connect to client: %v", err)
	}
	defer client.Close()

	res, err := client.Get("one")
	if err != nil {
		log.Println(err)
	}

	log.Println("one: ", res)

	// entry := map[string]interface{}{
	// 	"name": "John Doe",
	// 	"year": "2025",
	// }

	// response, err := client.Store("one", entry)
	// if err != nil {
	// 	log.Printf("error storing batch: %v", err)
	// }
	// log.Printf("stored batch, response: %v\n", response)

}

func Connect1() {
	client, err := tempdb.NewClient(tempdb.Config{
		Addr: "0.0.0.0:8081",
		URL:  "tempdb://admin:123456789@workspace:cb4552273c5c/ecommerce",
	})
	if err != nil {
		log.Fatalf("failed to connect to client: %v", err)
	}
	defer client.Close()

	ticker := time.NewTicker(3 * time.Second)
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

		product := make(map[string]interface{})
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
		log.Printf("stored product with key %s, response: %v\n", key, response)
	}
}
