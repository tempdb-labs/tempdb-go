# Go

This is a Go client library for interacting with the TempDB server.

## Installation

To use this client in your Go project, run:

```sh
go get github.com/ThembinkosiThemba/tempdb-go/lib
```

## Usage

Here's a basic example of how to use the TempDB Go client:

```go
package main

import (
	"fmt"
	"log"

	tempdb "github.com/ThembinkosiThemba/tempdb-go/lib"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	/// The NewCLient function initialises the client and takes 3 parameters
	/// 1. The address, either locally is ran there, or on a hosted client which is comming soon.
	/// 2. The database, this is the database you will be using to store data using the client.
	/// 3. Token, for access control, you will need to provide
	client, err := tempdb.NewClient("db-server-url", "ecommerce-store", apiKey)
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
}

```

Open this [file](test.go) for examples.

All methods return an error as the last return value. Always check this error to ensure your operations were successful.
