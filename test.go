package main

import (
	"fmt"
	"log"

	tempdb "github.com/ThembinkosiThemba/tempdb-go/lib"
)

func main() {
	client, err := tempdb.NewClient(tempdb.Config{
		Addr:       "0.0.0.0:8080",
		Collection: "tempdb://admin:admin@workspace:8020/ecommerce",
	})
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}
	defer client.Close()

	// Example usage when storing product information
	response, err := client.Store("productX", map[string]interface{}{
		"name":      "Laptop",
		"price":     999.99,
		"stock":     50,
		"Locations": "US",
	})
	if err != nil {
		log.Printf("Error setting user info: %v", err)
	} else {
		fmt.Println("Set user info response:", response)

	}

	// getting a particular product information
	getProductInfo, err := client.GetByKey("productX")
	if err != nil {
		log.Println("failed to get :", err)
	} else {
		log.Println("data: ", getProductInfo)
	}

	// data, err := client.ViewData()
	// if err != nil {
	// 	return
	// }

	// log.Println("data: ", data)

}
