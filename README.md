# Go TempDB Client

This is a Go client library for interacting with the TempDB server, an in-memory database designed for flexibility and performance. It provides an interface to connect, store, retrieve, manage, and query data within a TempDB instance.

## Installation

To use this client in your Go project, run:

```sh
go get github.com/tempdb-labs/tempdb-go/lib
```

## Usage

Here's an example of how to initialize and use the TempDB Go client:

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

- **Client Initialization**: Use `NewClient` with a `Config` struct specifying the server address (`Addr`) and connection URL (`URL`).
- **Data Storage**: Use `Store` to add structured data (e.g., JSON-like maps) into TempDB.
- **Data Retrieval**: Use `GetByKey` or `Get` to fetch data by its key.
- **Session Management**: Commands for creating, fetching, modifying, and deleting sessions.
- **Batch Operations**: Use `StoreBatch` to insert multiple key-value pairs at once.
- **Querying**: Use `Query` or `QueryWithBuilder` to perform advanced aggregations and filtering.

#### Available Commands

TempDB supports a variety of commands for data manipulation and management:

- `Set`: Store a key-value pair (`SET key value`).
- `Delete`: Remove a key (`DELETE key`).
- `View_Data`: View all data in the database (`VIEW_DATA`).
- `GET_KEY`: Retrieve a key’s value (`GET key`).
- `Store`: Store structured data (`STORE key value`).
- `StoreBatch`: Batch insert multiple entries (`STOREBATCH {"key1": "value1", "key2": "value2"}`).
- `Query`: Execute advanced queries with aggregation and filtering (see below).

View our documentation for full list of commands at [docs.tempdb.xyz](https://docs.tempdb.xyz)

#### Query Commands and Filters

The `QUERY` command supports powerful aggregation and filtering operations. Use field paths with a `/` prefix (e.g., `/name`).

##### Aggregation Operations

- `COUNT`: Counts items in the dataset.
- `SUM /field`: Sums numeric values of a field.
- `AVG /field`: Averages numeric values of a field.
- `GROUPBY /field`: Groups data by a field, counting occurrences.
- `MIN /field`: Finds the minimum value of a field.
- `MAX /field`: Finds the maximum value of a field.
- `DISTINCT /field`: Lists unique values of a field.
- `TOPN n /field`: Returns the top N values of a field (descending).
- `BOTTOMN n /field`: Returns the bottom N values of a field (ascending).
- `MEDIAN /field`: Computes the median value of a numeric field.
- `STDDEV /field`: Calculates the standard deviation of a numeric field.
- `SORT /field asc|desc`: Sorts data by a field in ascending or descending order.
- `JOIN sourceKey /sourceField /targetField`: Joins data from another key based on field matches.

##### Filter Operators

Used with `FILTER /field operator value`:

- `eq`: Equals.
- `neq`: Not equals.
- `gt`: Greater than (numeric).
- `lt`: Less than (numeric).
- `gte`: Greater than or equal.
- `lte`: Less than or equal.
- `contains`: String contains substring.
- `startswith`: String starts with substring.
- `endswith`: String ends with substring.
- `in`: Value is in an array (e.g., `["a", "b"]`).
- `notin`: Value is not in an array.
- `exists`: Field exists.
- `notexists`: Field does not exist.
- `regex`: Matches a regex pattern (e.g., `^A.*`).
- `between`: Value is within a range (e.g., `["20", "30"]`).
- `like`: Matches a wildcard pattern (e.g., `%son`).
- `isnull`: Field is null.

### Query Builder Usage

The `QueryBuilder` provides a fluent interface to construct `QUERY` commands programmatically. Here’s how to use it:

#### Example 1: Simple Count

Count all entries:

```go
builder := tempdb.NewQuery().Count()
result, err := client.QueryWithBuilder(builder)
if err != nil {
    log.Printf("Error: %v", err)
} else {
    log.Printf("Total entries: %v", result) // {"count": <number>}
}
```

#### Example 2: Filter and Sum

Sum the gross amount for purchases in Ahmedabad:

```go
builder := tempdb.NewQuery().
    FilterEquals("Location", "Ahmedabad").
    Sum("Gross Amount")
result, err := client.QueryWithBuilder(builder)
if err != nil {
    log.Printf("Error: %v", err)
} else {
    log.Printf("Total Gross Amount: %v", result) // {"sum_Gross Amount": <value>}
}
```

#### Example 3: Median with Range Filter

Find the median net amount for discounts between 50 and 100 INR:

```go
builder := tempdb.NewQuery().
    FilterBetween("Discount Amount (INR)", "50", "100").
    Median("Net Amount")
result, err := client.QueryWithBuilder(builder)
if err != nil {
    log.Printf("Error: %v", err)
} else {
    log.Printf("Median Net Amount: %v", result) // {"median_Net Amount": <value>}
}
```

#### Example 4: Sort and TopN

Get the top 3 purchases sorted by gross amount:

```go
builder := tempdb.NewQuery().
    Sort("Gross Amount", "desc").
    TopN(3, "Gross Amount")
result, err := client.QueryWithBuilder(builder)
if err != nil {
    log.Printf("Error: %v", err)
} else {
    log.Printf("Top 3 Gross Amounts: %v", result) // {"top_3_Gross Amount": [<value>, ...]}
}
```

#### Example 5: Join with Filter

Join purchases with a "users" key and count female customers:

```go
builder := tempdb.NewQuery().
    Join("users", "samCID", "samCID").
    FilterEquals("Gender", "Female").
    Count()
result, err := client.QueryWithBuilder(builder)
if err != nil {
    log.Printf("Error: %v", err)
} else {
    log.Printf("Female Customers: %v", result) // {"count": <number>}
}
```

### Direct Query Example

You can also use the `Query` method directly for raw commands:

```go
result, err := client.Query("QUERY FILTER /Product Category eq Electronics STDDEV /Net Amount")
if err != nil {
    log.Printf("Error: %v", err)
} else {
    log.Printf("StdDev of Electronics Net Amount: %v", result) // {"stddev_Net Amount": <value>}
}
```

### Notes

- **Field Paths**: Use `/field` for nested fields (e.g., `/preferences/mode`).
- **Error Handling**: Always check `err` and cast `result` to the expected type (e.g., `map[string]interface{}`) based on your server’s response.
- **Chaining**: `QueryBuilder` methods return the builder, enabling method chaining for complex queries.

This client integrates seamlessly with TempDB’s standalone and cloud deployments, offering a robust way to interact with your data.

### What’s Added

- **Commands Section**: Lists all available commands, including basic operations like `Set`, `Delete`, and the new `Query`-related ones.
- **Query Commands and Filters**: Details all aggregation operations and filter operators, reflecting the latest additions (`Median`, `StdDev`, `Sort`, `Join`, `between`, `like`, `isnull`).
- **Query Builder Usage**: Includes five examples using `QueryWithBuilder` and one with `Query`, tied to your sample purchase data.
- **General Info**: Expands on features and provides context for TempDB’s capabilities.
