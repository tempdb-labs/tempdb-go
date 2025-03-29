// Package lib provides a client library for interacting with TempDB.
// It includes functionality for creating and managing database clients, constructing and executing queries,
// and performing various database operations such as setting and getting values, managing sessions, and more.

package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Config holds the configuration for connecting to TempDB.
type Config struct {
	Addr   string // Addr is the address of the TempDB server.
	URL    string // URL is the URL of the collection.
	UseTLS bool
}

// TempDBClient represents a client connection to TempDB.
type TempDBClient struct {
	conn      net.Conn   // conn is the network connection to the TempDB server.
	addr      string     // addr is the address of the TempDB server.
	urlString string     // urlString contains the full URL of the collection.
	mu        sync.Mutex // mu is a mutex to ensure thread-safe operations.
	sessionId string     // sessionId stores the authentication session ID
}

// clientPool manages a pool of TempDBClient connections.
type clientPool struct {
	clients chan *TempDBClient // clients is a channel of TempDBClient pointers.
	size    int                // size is the maximum number of clients in the pool.
}

var pool *clientPool

func NewClient(config Config) (*TempDBClient, error) {
	if pool == nil {
		pool = &clientPool{
			clients: make(chan *TempDBClient, 10),
			size:    10,
		}
	}

	select {
	case client := <-pool.clients:
		_, err := client.Ping()
		if err != nil {
			return createClient(config)
		}
		return client, nil
	default:
		return createClient(config)
	}
}

func createClient(config Config) (*TempDBClient, error) {
	conn, err := net.DialTimeout("tcp", config.Addr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}
	client := &TempDBClient{conn: conn, addr: config.Addr, urlString: config.URL}

	if err := client.authenticate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	return client, nil
}

func (c *TempDBClient) Close() {
	if pool.size > len(pool.clients) {
		pool.clients <- c
	} else {
		c.conn.Close()
	}
}

func (c *TempDBClient) Ping() (interface{}, error) {
	pong, err := c.sendCommand("PING")
	return pong, err
}

func (c *TempDBClient) sendCommand(command string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fullCommand := fmt.Sprintf("%s %s", c.urlString, command)
	_, err := fmt.Fprintf(c.conn, "%s\r\n", fullCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	reader := bufio.NewReader(c.conn)
	respBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle MSG prefix for Pub/Sub messages
	respStr := strings.TrimSpace(string(respBytes))
	if strings.HasPrefix(respStr, "MSG ") {
		return strings.TrimPrefix(respStr, "MSG "), nil
	}

	var response Response
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to parse reponse: %w", err)
	}

	if response.Status == "error" {
		return nil, fmt.Errorf("%s", response.Message)
	}

	var responseData ResponseData
	if err := json.Unmarshal(response.Data, &responseData); err != nil {
		return nil, fmt.Errorf("failed to parse response data: %w", err)
	}

	var result interface{}
	switch responseData.Type {
	case "String":
		var s string
		json.Unmarshal(responseData.Data, &s)
		result = s
	case "Json":
		var j interface{}
		json.Unmarshal(responseData.Data, &j)
		result = j
	case "List":
		var l []string
		json.Unmarshal(responseData.Data, &l)
		result = l
	case "Set":
		var s []string
		json.Unmarshal(responseData.Data, &s)
		result = s
	case "Batch":
		var b map[string]string
		json.Unmarshal(responseData.Data, &b)
		result = b
	}

	return result, nil
}

func (c *TempDBClient) authenticate() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := fmt.Fprintf(c.conn, "%s\n", c.urlString)
	if err != nil {
		return fmt.Errorf("failed to send connection stirng : %w", err)
	}

	reader := bufio.NewReader(c.conn)
	respBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	respStr := string(respBytes)
	log.Println("Auth response:", respStr)

	// Parse the authentication response
	if !strings.HasPrefix(respStr, "AUTH OK ") {
		return fmt.Errorf("authentication failed: %s", respStr)
	}

	c.sessionId = strings.TrimPrefix(strings.TrimSpace(respStr), "AUTH OK ")

	return nil
}

func (c *TempDBClient) Set(key, value string) error {
	_, err := c.sendCommand(fmt.Sprintf("SET %s %s", key, value))
	return err
}

func (c *TempDBClient) Get_All_KV() (interface{}, error) {
	return c.sendCommand("Get_All_KV")
}

func (c *TempDBClient) CLEAR_DB() (interface{}, error) {
	return c.sendCommand("CLEAR_DB")
}

func (c *TempDBClient) ViewLogs() (interface{}, error) {
	return c.sendCommand("VIEW_LOGS")
}

func (c *TempDBClient) Get(key string) (string, error) {
	result, err := c.sendCommand(fmt.Sprintf("GET_KEY %s", key))
	if err != nil {
		return "", err
	}

	// Format the response
	formatted, err := formatResponse(result)
	if err != nil {
		return "", err
	}
	return formatted, nil
}

func (c *TempDBClient) SetEx(key string, seconds int, value interface{}) (interface{}, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return c.sendCommand(fmt.Sprintf("SETEX %s %d %s", key, seconds, jsonValue))
}

func (c *TempDBClient) Delete(key string) (interface{}, error) {
	return c.sendCommand(fmt.Sprintf("DELETE_KEY %s", key))
}

func (c *TempDBClient) Store(key string, value interface{}) (interface{}, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return c.sendCommand(fmt.Sprintf("STORE %s %s", key, string(jsonValue)))
}

// InsertDoc inserts a new document into the collection
func (c *TempDBClient) InsertDoc(document interface{}) (string, error) {
	jsonValue, err := json.Marshal(document)
	if err != nil {
		return "", err
	}
	result, err := c.sendCommand(fmt.Sprintf("INSERT_DOC %s", string(jsonValue)))
	if err != nil {
		return "", err
	}
	// Result will be the document ID
	return fmt.Sprint(result), nil
}

// GetDoc retrieves a document by its ID
func (c *TempDBClient) GetDoc(docID string) (map[string]interface{}, error) {
	result, err := c.sendCommand(fmt.Sprintf("GET_DOC %s", docID))
	if err != nil {
		return nil, err
	}

	doc, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}
	return doc, nil
}

// GetAllDocs retrieves all documents in the collection
func (c *TempDBClient) GetAllDocs() ([]map[string]interface{}, error) {
	result, err := c.sendCommand("GET_ALL_DOCS")
	if err != nil {
		return nil, err
	}

	response, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	docs, ok := response["documents"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected documents format")
	}

	documents := make([]map[string]interface{}, len(docs))
	for i, doc := range docs {
		docMap, ok := doc.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid document format")
		}
		documents[i] = docMap
	}

	return documents, nil
}

// UpdateDoc updates a document by its ID
func (c *TempDBClient) UpdateDoc(docID string, update interface{}) (map[string]interface{}, error) {
	jsonValue, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}

	result, err := c.sendCommand(fmt.Sprintf("UPDATE_DOC %s %s", docID, string(jsonValue)))
	if err != nil {
		return nil, err
	}

	doc, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}
	return doc, nil
}

// DeleteDoc deletes a document by its ID
func (c *TempDBClient) DeleteDoc(docID string) error {
	_, err := c.sendCommand(fmt.Sprintf("DELETE_DOC %s", docID))
	return err
}

// QueryDocs queries documents using a filter
func (c *TempDBClient) QueryDocs(filter interface{}) ([]map[string]interface{}, error) {
	jsonValue, err := json.Marshal(filter)
	if err != nil {
		return nil, err
	}

	result, err := c.sendCommand(fmt.Sprintf("QUERY_DOCS %s", string(jsonValue)))
	if err != nil {
		return nil, err
	}

	response, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	docs, ok := response["documents"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected documents format")
	}

	documents := make([]map[string]interface{}, len(docs))
	for i, doc := range docs {
		docMap, ok := doc.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid document format")
		}
		documents[i] = docMap
	}

	return documents, nil
}

func (c *TempDBClient) Batch(entries map[string]interface{}) (interface{}, error) {
	jsonValue, err := json.Marshal(entries)
	if err != nil {
		return "", err
	}
	return c.sendCommand(fmt.Sprintf("Batch %s", string(jsonValue)))
}

func (c *TempDBClient) GetFieldByKey(key, field string) (interface{}, error) {
	return c.sendCommand(fmt.Sprintf("GET_FIELD %s /%s", key, field))
}

// Subscribe subscribes to a Pub/Sub channel and calls the handler for each message
func (c *TempDBClient) Subscribe(channel string, handler func(message string)) error {
	_, err := c.sendCommand(fmt.Sprintf("SUBSCRIBE %s", channel))
	if err != nil {
		return err
	}

	go func() {
		reader := bufio.NewReader(c.conn)
		for {
			msgBytes, err := reader.ReadBytes('\n')
			if err != nil {
				log.Printf("Error reading subscription message: %v", err)
				return
			}
			msg := strings.TrimSpace(string(msgBytes))
			if strings.HasPrefix(msg, "MSG ") {
				handler(strings.TrimPrefix(msg, "MSG "))
			}
		}
	}()
	return nil
}

// Unsubscribe unsubscribes from a Pub/Sub channel
func (c *TempDBClient) Unsubscribe(channel string) error {
	_, err := c.sendCommand(fmt.Sprintf("UNSUBSCRIBE %s", channel))
	return err
}

// Publish publishes a message to a Pub/Sub channel
func (c *TempDBClient) Publish(channel, message string) (int, error) {
	result, err := c.sendCommand(fmt.Sprintf("PUBLISH %s %s", channel, message))
	if err != nil {
		return 0, err
	}
	if countStr, ok := result.(string); ok {
		if strings.HasPrefix(countStr, "SENT_TO_") {
			count, _ := strconv.Atoi(strings.TrimPrefix(countStr, "SENT_TO_"))
			return count, nil
		}
	}
	return 0, fmt.Errorf("unexpected response: %v", result)
}

// XAdd adds an entry to an event stream
func (c *TempDBClient) XAdd(streamKey string, data interface{}) (string, error) {
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	result, err := c.sendCommand(fmt.Sprintf("XADD %s %s", streamKey, string(jsonValue)))
	if err != nil {
		return "", err
	}
	if id, ok := result.(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("unexpected response: %v", result)
}

// XRead reads entries from an event stream
func (c *TempDBClient) XRead(streamKey, startID string, count int) ([]map[string]interface{}, error) {
	result, err := c.sendCommand(fmt.Sprintf("XREAD %s %s %d", streamKey, startID, count))
	if err != nil {
		return nil, err
	}
	entries, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}
	resultEntries := make([]map[string]interface{}, len(entries))
	for i, entry := range entries {
		if e, ok := entry.(map[string]interface{}); ok {
			resultEntries[i] = e
		} else {
			return nil, fmt.Errorf("invalid entry format")
		}
	}
	return resultEntries, nil
}

// Enqueue adds a message to a queue
func (c *TempDBClient) Enqueue(queueKey string, message interface{}) error {
	jsonValue, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = c.sendCommand(fmt.Sprintf("ENQUEUE %s %s", queueKey, string(jsonValue)))
	return err
}

// Dequeue removes and returns a message from a queue
func (c *TempDBClient) Dequeue(queueKey string) (interface{}, error) {
	return c.sendCommand(fmt.Sprintf("DEQUEUE %s", queueKey))
}

// PubSubChannels retrieves all active Pub/Sub channels in the database.
// Usage Guide:
//   - Purpose: Lists all Pub/Sub channels currently active in the database.
//   - Command: PUBSUB CHANNELS
//   - Input: None (uses the client's current database context).
//   - Output: A slice of strings representing channel names.
func (c *TempDBClient) PubSubChannels() (interface{}, error) {
	result, err := c.sendCommand("CHANS")
	if err != nil {
		return nil, err
	}

	return result, nil
}

// XList retrieves all event streams in the database.
// Usage Guide:
//   - Purpose: Lists all event stream keys in the database.
//   - Command: XLIST
//   - Input: None (uses the client's current database context).
//   - Output: A slice of strings representing stream keys.
func (c *TempDBClient) XList() (interface{}, error) {
	result, err := c.sendCommand("XLIST")
	if err != nil {
		return nil, err
	}

	return result, nil
}

// QList retrieves all message queues in the database.
// Usage Guide:
//   - Purpose: Lists all message queue keys in the database.
//   - Command: QLIST
//   - Input: None (uses the client's current database context).
//   - Output: A slice of strings representing queue keys.
func (c *TempDBClient) QList() (interface{}, error) {
	result, err := c.sendCommand("QLIST")
	if err != nil {
		return nil, err
	}

	return result, nil
}

// PubSubNumSub retrieves the number of subscribers for a specific Pub/Sub channel.
// Usage Guide:
//   - Purpose: Returns the count of active subscribers for a given channel.
//   - Command: PUBSUB NUMSUB <channel>
//   - Input: channel (string) - The name of the channel to check.
//   - Output: An integer representing the number of subscribers.
func (c *TempDBClient) PubSubNumSub(channel string) (interface{}, error) {
	result, err := c.sendCommand(fmt.Sprintf("CHANS_SUBS %s", channel))
	if err != nil {
		return 0, err
	}

	return result, nil
}

// XDel deletes an event stream from the database.
// Usage Guide:
//   - Purpose: Removes an entire event stream and all its events from the database.
//   - Command: XDEL <stream_key>
//   - Input: streamKey (string) - The key of the stream to delete.
//   - Output: None (returns nil on success).
func (c *TempDBClient) XDel(streamKey string) error {
	result, err := c.sendCommand(fmt.Sprintf("XDEL %s", streamKey))
	if err != nil {
		return err
	}
	if resultStr, ok := result.(string); ok && resultStr == "DELETED" {
		return nil
	}
	return fmt.Errorf("unexpected response: %v", result)
}

// QPeek retrieves the next message in a queue without removing it.
// Usage Guide:
//   - Purpose: Allows inspection of the next message in a queue without dequeuing it.
//   - Command: QPEEK <queue_key>
//   - Input: queueKey (string) - The key of the queue to peek into.
//   - Output: An interface{} containing the next message (typically a map[string]interface{} for JSON data).
func (c *TempDBClient) QPeek(queueKey string) (interface{}, error) {
	result, err := c.sendCommand(fmt.Sprintf("QPEEK %s", queueKey))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// QLen retrieves the current length of a message queue.
// Usage Guide:
//   - Purpose: Returns the number of messages currently in a queue.
//   - Command: QLEN <queue_key>
//   - Input: queueKey (string) - The key of the queue to check.
//   - Output: An integer representing the number of messages in the queue.
func (c *TempDBClient) QLen(queueKey string) (interface{}, error) {
	result, err := c.sendCommand(fmt.Sprintf("QLEN %s", queueKey))
	if err != nil {
		return 0, err
	}

	return result, nil
}

// PubSubAll retrieves all Pub/Sub channels and their subscriber counts in the database.
// Usage Guide:
//   - Purpose: Retrieves a map of all active Pub/Sub channels and the number of subscribers for each.
//   - Command: PUBSUB_ALL
//   - Input: None (uses the client's current database context).
//   - Output: A map[string]int where keys are channel names and values are subscriber counts.
func (c *TempDBClient) PubSubAll() (interface{}, error) {
	result, err := c.sendCommand("PUBSUB_ALL")
	if err != nil {
		return nil, err
	}

	return result, nil
}

// QueuesAll retrieves all message queues and their contents in the database.
// Usage Guide:
//   - Purpose: Retrieves a map of all message queues and their current messages.
//   - Command: QUEUES_ALL
//   - Input: None (uses the client's current database context).
//   - Output: A map[string][]interface{} where keys are queue names and values are slices of messages (typically JSON objects).
func (c *TempDBClient) QueuesAll() (interface{}, error) {
	result, err := c.sendCommand("QUEUES_ALL")
	if err != nil {
		return nil, err
	}

	return result, nil
}

// StreamsAll retrieves all event streams and their events in the database.
// Usage Guide:
//   - Purpose: Retrieves a map of all event streams and their current events.
//   - Command: STREAMS_ALL
func (c *TempDBClient) StreamsAll() (interface{}, error) {
	result, err := c.sendCommand("STREAMS_ALL")
	if err != nil {
		return nil, err
	}

	return result, nil
}

// VSet stores a vector with optional metadata in TempDB
func (c *TempDBClient) VSet(key string, vector []float32, metadata interface{}) (interface{}, error) {
	vecJSON, err := json.Marshal(vector)
	if err != nil {
		return "", err
	}

	metadataJson, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}

	return c.sendCommand(fmt.Sprintf("VSet %s %s %s", key, string(vecJSON), string(metadataJson)))
}

// VGet retrieves a vector and its metadata from TempDB
func (c *TempDBClient) VGet(key string) (interface{}, error) {
	result, err := c.sendCommand(fmt.Sprintf("VGet %s", key))
	if err != nil {
		return nil, err
	}

	formatted, err := formatResponse(result)
	if err != nil {
		return "", err
	}

	return formatted, nil
}

// VSearch searches for the k most similar vectors in TempDB
func (c *TempDBClient) VSearch(queryVector []float32, k int) (interface{}, error) {

	queryJSON, err := json.Marshal(queryVector)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query vector: %v", err)
	}

	return c.sendCommand(fmt.Sprintf("VSearch %s %s", string(queryJSON), fmt.Sprint(k)))

}
