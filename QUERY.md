# TempDB Query Documentation

## Overview
TempDB now supports two querying approaches:
1. **Query Builder**: A structured, type-safe way to build queries using method chaining
2. **Pipeline**: A string-based syntax for creating queries using a pipeline pattern

## Query Builder Approach

### Basic Usage
```go
client, _ := NewClient(Config{
    Addr: "localhost:8080",
    Collection: "mycollection",
})

query := NewQueryBuilder().
    WhereEq("status", "active").
    Sort("created_at", Descending).
    Limit(10).
    Build()

result, err := client.Query(query)
```

### Available Methods

#### Conditions
- `WhereEq(field, value string)`: Exact match condition
```go
queryBuilder.WhereEq("status", "active")
```

#### Sorting
- `Sort(field string, order SortOrder)`: Sort results by field
```go
queryBuilder.Sort("created_at", Descending)
// or
queryBuilder.Sort("age", Ascending)
```

#### Pagination
- `Limit(limit int)`: Limit number of results
- `Offset(offset int)`: Skip first N results
```go
queryBuilder.Limit(10).Offset(20) // Get results 21-30
```

#### Time Range
- `TimeRange(start, end int64)`: Filter by creation timestamp
```go
queryBuilder.TimeRange(startTime, endTime)
```

### Complete Example
```go
result, err := client.Query(
    NewQueryBuilder().
        WhereEq("status", "active").
        WhereEq("type", "user").
        Sort("created_at", Descending).
        Limit(10).
        Offset(0).
        TimeRange(startTime, endTime).
        Build(),
)
```

## Pipeline Approach

### Basic Usage
```go
pipeline := "where status eq active | sort created_at desc | limit 10"
result, err := client.QueryPipeline(pipeline)
```

### Pipeline Operators

#### Where Clause
```
where <field> <operator> <value>
```
Supported operators:
- `eq`: Equal to
- `gt`: Greater than
- `lt`: Less than

Example:
```go
"where age gt 25"
"where status eq active"
```

#### Sort
```
sort <field> <order>
```
Order options:
- `asc`: Ascending
- `desc`: Descending

Example:
```go
"sort created_at desc"
"sort name asc"
```

#### Limit
```
limit <number>
```
Example:
```go
"limit 10"
```

### Chaining Pipeline Operations
Use the pipe character (`|`) to chain operations:
```go
pipeline := "where status eq active | sort created_at desc | limit 10"
result, err := client.QueryPipeline(pipeline)
```

### Pipeline Examples

1. Find active users, sorted by creation date:
```go
"where status eq active | sort created_at desc"
```

2. Get the 10 newest records:
```go
"sort created_at desc | limit 10"
```

3. Filter and limit results:
```go
"where type eq user | where age gt 25 | limit 5"
```

## Working with Results

Both query approaches return results in JSON format. The response will be a string containing an array of JSON objects that match the query criteria.

Example response:
```json
[
  {
    "id": "1",
    "status": "active",
    "created_at": 1641234567,
    "name": "John Doe"
  },
  {
    "id": "2",
    "status": "active",
    "created_at": 1641234568,
    "name": "Jane Smith"
  }
]
```

## Error Handling

Both query methods return an error as their second return value:
```go
result, err := client.Query(query)
if err != nil {
    // Handle error
}
```

Common error cases:
- Invalid query syntax
- Invalid field names
- Invalid operator usage
- Network errors
- Server-side processing errors

## Performance Considerations

1. **Indexing**: The database automatically creates indexes based on frequently accessed fields
2. **Query Complexity**: Pipeline operations are processed sequentially, so order them for optimal performance (e.g., filter before sort)
3. **Result Size**: Use limits to restrict the result set size when dealing with large datasets
4. **Time Range Queries**: Using TimeRange in Query Builder or combining with other filters can help reduce the result set

## Best Practices

1. **Use Query Builder for Complex Queries**
   - Provides type safety
   - Better IDE support
   - Easier to maintain

2. **Use Pipeline for Simple Queries**
   - More concise for simple operations
   - Easier to construct dynamically
   - Familiar syntax for SQL users

3. **Pagination**
   - Always use limits when displaying results
   - Combine offset and limit for pagination
   - Consider using sort with pagination for consistent results

4. **Error Handling**
   - Always check for errors in returned results
   - Implement proper error handling for network issues
   - Log query errors for debugging

5. **Resource Management**
   ```go
   client, err := NewClient(config)
   if err != nil {
       // Handle error
   }
   defer client.Close()
   ```

## Limitations

Current implementation limitations:
1. No support for complex boolean operations (AND, OR, NOT)
2. Limited set of comparison operators
3. No support for aggregations or grouping
4. No support for nested field queries
5. No support for regex or pattern matching