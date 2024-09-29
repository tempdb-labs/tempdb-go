package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	tempdb "github.com/ThembinkosiThemba/tempdb-libs/go"
// )

// func main() {
// 	apiKey := os.Getenv("API_KEY")
// 	client, err := tempdb.NewClient("db-server-url", "testing", apiKey)
// 	if err != nil {
// 		log.Fatalf("Failed to get client: %v", err)
// 	}
// 	defer client.Close()

// 	// Example usage when storing product information
// 	response, err := client.Store("productX", map[string]interface{}{
// 		"name":      "Laptop",
// 		"price":     999.99,
// 		"stock":     50,
// 		"Locations": "US",
// 	})
// 	if err != nil {
// 		log.Printf("Error setting user info: %v", err)
// 	} else {
// 		fmt.Println("Set user info response:", response)

// 	}

// 	// getting a particular product information
// 	getProductInfo, err := client.GetByKey("productX")
// 	if err != nil {
// 		log.Println("failed to get :", err)
// 	} else {
// 		log.Println("data: ", getProductInfo)
// 	}
// }
