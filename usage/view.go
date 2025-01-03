package main

import (
	"log"

	tempdb "github.com/tempdb-labs/tempdb-go/lib"
)

func main() {
	client, err := tempdb.NewClient(tempdb.Config{
		Addr: "0.0.0.0:8081",
		URL:  "tempdb://admin:admin@workspace:8020/ecommerce",
	})
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}
	defer client.Close()

	res, err := client.ViewData()
	if err != nil {
		return
	}

	log.Println("data: ", res)
}
