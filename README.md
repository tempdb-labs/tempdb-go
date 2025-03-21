## Go TempDB Client

The Go TempDB Client is a powerful, lightweight library for interacting with TempDB, an in-memory database optimized for speed, flexibility, and ease of use. TempDB supports multiple data paradigms including key-value storage, document storage, vector operations, and messaging (Pub/Sub, queues, and event streams), making it an ideal choice for real-time applications, analytics, and dynamic workloads.

This client provides a simple Go interface to connect to a TempDB server, manage data, and execute queries with a `QueryBuilder` or `raw command strings`.

**Key Features:**

- Multi-paradigm support: Key-value, documents, vectors, and messaging.
- Thread-safe client pooling for efficient connections.
- Advanced querying with aggregation, filtering, and joining capabilities.

### Table of Contents

- [Quick Start](#quick-start)
- [Commands](#commands)
  - [Key-Value Commands](#key-value-commands)
  - [Document Commands](#document-commands)
  - [Vector Commands](#vector-commands)
  - [Messaging Commands](#messaging-commands)
- [Querying the Database](#querying-the-database)
  - [Using Raw Query Strings](#using-raw-query-strings)
  - [Using QueryBuilder](#using-querybuilder)
  - [All Query Commands and Filters](#query-commands-and-filters)
- [More Information](#more-information)

### Quick Start

Get up and running with the Go TempDB Client in minutes. This section demonstrates how to install the library and perform a basic operation.

#### Installation

Add the TempDB client to your Go project:

```sh
go get github.com/tempdb-labs/tempdb-go/lib
```

#### Example: Storing and Retrieving User Data

Below is a simple example of connecting to TempDB, storing a user object, and retrieving it:

Here's an example of how to initialize and use the TempDB Go client:

```go
package main

import (
	"log"
	tempdb "github.com/tempdb-labs/tempdb-go/lib"
)

func main() {
	// Initialize the client
	client, err := tempdb.NewClient(tempdb.Config{
		Addr: "0.0.0.0:8081",
		URL:  "tempdb://admin:admin@users:kv",
	})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	// Define a user object
	user := map[string]interface{}{
		"name":    "Thembinkosi",
		"surname": "Mkhonta",
		"preferences": map[string]interface{}{
			"mode":          "dark",
			"notifications": "no",
		},
	}

	// Store the user
	result, err := client.Store("user_01", user)
	if err != nil {
		log.Fatalf("Failed to store user: %v", err)
	}
	log.Printf("User stored: %v", result)

	// Retrieve the user
	retrieved, err := client.Get("user_01")
	if err != nil {
		log.Fatalf("Failed to retrieve user: %v", err)
	}
	log.Printf("Retrieved user: %s", retrieved)
}
```

**Notes**

- Replace `Addr` and `URL` with your TempDB server’s configuration.
- The URL format is `tempdb://<username>:<password>@<database>:<database_type>`.
- Always defer client.Close() to return the connection to the pool or close it.

### Commands

The TempDB Go Client supports a variety of commands organized by data paradigm. Below is a comprehensive list with their respective methods and use cases.

#### Key-Value Commands

For basic key-value storage and retrieval:

- **`Set(key, value string) error`**: Stores a simple key-value pair.
  - Example: `client.Set("name", "Alice")`
- **`Get(key string) (string, error)`**: Retrieves a value by key.
  - Example: `value, err := client.Get("name")`
- **`Delete(key string) (interface{}, error)`**: Deletes a key-value pair.
  - Example: `client.Delete("name")`
- **`Store(key string, value interface{}) (interface{}, error)`**: Stores structured data (e.g., JSON-like maps).
  - Example: `client.Store("user_01", map[string]interface{}{"name": "Bob"})`
- **`SetEx(key string, seconds int, value interface{}) (interface{}, error)`**: Stores data with an expiration time.
  - Example: `client.SetEx("temp_key", 60, "temp_value")`
- **`Get_All_KV() (interface{}, error)`**: Returns all associated data with your key-value database.
  - Example: `client.Get_All_KV()`
- **`Batch(entries map[string]interface{}) (interface{}, error)`**: Stores multiple key-value pairs in one operation.
  - Example: `client.Batch(map[string]interface{}{"k1": "v1", "k2": "v2"})`

#### Document Commands

For document-oriented storage:

- **`InsertDoc(document interface{}) (string, error)`**: Inserts a document and returns its ID.
  - Example: `docID, err := client.InsertDoc(map[string]interface{}{"name": "John"})`
- **`GetDoc(docID string) (map[string]interface{}, error)`**: Retrieves a document by ID.
  - Example: `doc, err := client.GetDoc("doc123")`
- **`GetAllDocs() ([]map[string]interface{}, error)`**: Retrieves all documents in the collection.
  - Example: `docs, err := client.GetAllDocs()`
- **`UpdateDoc(docID string, update interface{}) (map[string]interface{}, error)`**: Updates a document by ID.
  - Example: `updated, err := client.UpdateDoc("doc123", map[string]interface{}{"age": 31})`
- **`DeleteDoc(docID string) error`**: Deletes a document by ID.
  - Example: `client.DeleteDoc("doc123")`
- **`QueryDocs(filter interface{}) ([]map[string]interface{}, error)`**: Queries documents with a filter.
  - Example: `docs, err := client.QueryDocs(map[string]interface{}{"age": 30})`

#### Vector Commands

For vector storage and similarity search:

- **`VSet(key string, vector []float32, metadata interface{}) (interface{}, error)`**: Stores a vector with metadata.
  - Example: `client.VSet("vec1", []float32{1.0, 2.0}, map[string]string{"type": "image"})`
- **`VGet(key string) (interface{}, error)`**: Retrieves a vector and its metadata.
  - Example: `vec, err := client.VGet("vec1")`
- **`VSearch(queryVector []float32, k int) (interface{}, error)`**: Finds the k most similar vectors.
  - Example: `results, err := client.VSearch([]float32{1.1, 2.1}, 5)`

#### Messaging Commands

For Pub/Sub, queues, and event streams:

- **Pub/Sub:**

  - **`Subscribe(channel string, handler func(message string)) error`**: Subscribes to a channel.
    - Example: `client.Subscribe("news", func(msg string) { log.Println(msg) })`
  - **`Unsubscribe(channel string) error`**: Unsubscribes from a channel.
    - Example: `client.Unsubscribe("news")`
  - **`Publish(channel, message string) (int, error)`**: Publishes a message to a channel.
    - Example: `count, err := client.Publish("news", "Hello")`
  - **`PubSubChannels() (interface{}, error)`**: Lists all active channels.
    - Example: `channels, err := client.PubSubChannels()`

- **Event Streams:**

  - **`XAdd(streamKey string, data interface{}) (string, error)`**: Adds an entry to a stream.
    - Example: `id, err := client.XAdd("orders", map[string]string{"item": "phone"})`
  - **`XRead(streamKey, startID string, count int) ([]map[string]interface{}, error)`**: Reads stream entries.
    - Example: `entries, err := client.XRead("orders", "0", 10)`
  - **`XDel(streamKey string) error`**: Deletes a stream.
    - Example: `client.XDel("orders")`

- **Queues:**
  - **`Enqueue(queueKey string, message interface{}) error`**: Adds a message to a queue.
    - Example: `client.Enqueue("tasks", "process_order")`
  - **`Dequeue(queueKey string) (interface{}, error)`**: Removes and returns a message.
    - Example: `msg, err := client.Dequeue("tasks")`
  - **`QPeek(queueKey string) (interface{}, error)`**: Peeks at the next message without removing it.
    - Example: `msg, err := client.QPeek("tasks")`
  - **`QLen(queueKey string) (interface{}, error)`**: Returns the queue length.
    - Example: `length, err := client.QLen("tasks")`

#### Common commands

- **`CLEAR_DB() (interface{}, error)`**: Clears and drops the database.
  - Example: `clearDb, err := client.CLEAR_DB()`
- **`VIEW_LOGS() (interface{}, error)`**: Views you own database logs.
  - Example: `logs, err := client.VIEW_LOGS()`

VIEW_LOGS

### Querying the Database

TempDB supports querying capabilities through the `Query` method (raw strings) or the `QueryBuilder` (fluent API). This section explains both approaches and provides e-commerce case studies.

#### Using Raw Query Strings

The `Query` method accepts a raw command string in the format: `QUERY <operations>`. Use field paths with a `/` prefix (e.g., `/name`).

**Syntax**

- **Aggregations**: `COUNT`, `SUM /field`, `AVG /field`, `GROUPBY /field`, etc.
- **Filters**: `FILTER /field operator value` (e.g., `FILTER /age gt 25`).
- **Order**: Combine operations with spaces (e.g., `FILTER /age gt 25 GROUPBY /city COUNT`).

##### Example: Total Sales by Category

```go
result, err := client.Query("GROUPBY /category SUM /net_amount")
if err != nil {
    log.Printf("Error: %v", err)
} else {
    log.Printf("Sales by category: %v", result) // e.g., {"category": {"Electronics": 5000, "Clothing": 3000}}
}

```

#### Using QueryBuilder

The `QueryBuilder` provides a programmatic way to construct queries with `method chaining`. It’s ideal for complex logic and maintainability.

**Syntax**

- Start with `NewQuery()`.
- Chain methods like `.Filter(), .GroupBy(), .Sum()`, etc.
- Finalize with `.Build()` or pass directly to `QueryWithBuilder()`.

  **Example: Average Purchase by Gender**

```go
builder := tempdb.NewQuery().
    GroupBy("gender").
    Average("net_amount")
result, err := client.QueryWithBuilder(builder)
if err != nil {
    log.Printf("Error: %v", err)
} else {
    log.Printf("Average by gender: %v", result) // e.g., {"avg_net_amount": {"Male": 1500, "Female": 1800}}
}
```

#### Query Commands and Filters

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

### NB

- **Field Paths**: Use `/field` for nested fields (e.g., `/preferences/mode`).
- **Error Handling**: Always check `err` and cast `result` to the expected type (e.g., `map[string]interface{}`) based on your server’s response.
- **Chaining**: `QueryBuilder` methods return the builder, enabling method chaining for complex queries.

### More Information

For detailed documentation, additional examples, and TempDB server setup instructions, visit:

- **Official Documentation**: [docs.tempdb.xyz](https://docs.tempdb.xyz)
- **GitHub Repository**: [github.com/tempdb-labs/tempdb-go](https://github.com/tempdb-labs/tempdb-go)
- **Support**: Reach out via GitHub issues or email `mkhonta@tempdb.xyz`.

Stay updated with the latest features and contribute to the project on GitHub!
