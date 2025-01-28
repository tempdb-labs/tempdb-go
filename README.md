# Go

This is a Go client library for interacting with the TempDB server. It provides an interface to connect, store, retrieve, and manage data within a TempDB database.

## Installation

To use this client in your Go project, run:

```sh
go get github.com/tempdb-labs/tempdb-go/lib
```

## Usage

Here's an example of how to use the TempDB Go client:

```go
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

	user := map[string]interface{}{
		"name":    "Thembinkosi",
		"surname": "Mkhonta",
		"preferences": map[string]interface{}{
			"mode":          "dark",
			"notifications": "no",
		},
	}

	result, err := client.Store("user_01", user)
	if err != nil {
		log.Printf("Failed to store user: %v", err)
		return
	}

	log.Printf("User stored successfully: %v", result)
}
```

### Key Features and Functions

- **Client Initialization**: Create a new client using `NewClient` with `Config` containing the server address and collection URL.
- **Data Storage**: Use `Store` to add structured data into TempDB.
- **Data Retrieval**: Use `GetByKey` to fetch data by its key.
- **Session Management**: Includes commands for creating, fetching, modifying, and deleting sessions.
- **Batch Operations**: Use `StoreBatch` to insert multiple entries at once.
- **Other Commands**:
  - `Set`: Store a key-value pair.
  - `SetEx`: Store a key-value pair with an expiry.
  - `Delete`: Remove a key from the database.
  - `LPush`: Push a value onto a list.
  - `SAdd`: Add a value to a set.
  - `GetFieldByKey`: Retrieve specific fields within a key.
  - `ViewData`: View all data in the database.
  - `Get`: Retrieves a particular key from the database.

### Example for Data Retrieval

```go
getProductInfo, err := client.Get("user_01")
if err != nil {
	log.Printf("Failed to get product: %v", err)
} else {
	log.Printf("Product data: %v", getProductInfo)
}
```

### Session Management Example

```go
sessionID, err := client.CreateSession("user_01")
if err != nil {
	log.Printf("Failed to create session: %v", err)
}
log.Printf("Session created: %v", sessionID)
```

Refer to the [library full documentation](https://docs.tempdb.xyz) for further details on all available commands and their usage.
