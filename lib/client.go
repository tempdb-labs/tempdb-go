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
	"sync"
	"time"
)

// Config holds the configuration for connecting to TempDB.
type Config struct {
	Addr string // Addr is the address of the TempDB server.
	URL  string // URL is the URL of the collection.
}

// TempDBClient represents a client connection to TempDB.
type TempDBClient struct {
	conn      net.Conn   // conn is the network connection to the TempDB server.
	addr      string     // addr is the address of the TempDB server.
	urlString string     // urlString contains the full URL of the collection.
	mu        sync.Mutex // mu is a mutex to ensure thread-safe operations.
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
	return &TempDBClient{conn: conn, addr: config.Addr, urlString: config.URL}, nil
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

func Client(addr, url string) (*TempDBClient, error) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	client := &TempDBClient{conn: conn, addr: addr, urlString: url}

	return client, nil
}

func (c *TempDBClient) sendCommand(command string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fullCommand := fmt.Sprintf("%s %s", c.urlString, command)
	_, err := fmt.Fprintf(c.conn, "%s", fullCommand+"\r\n")
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(c.conn)
	respBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
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

func (c *TempDBClient) Set(key, value string) error {
	_, err := c.sendCommand(fmt.Sprintf("SET %s %s", key, value))
	return err
}

func (c *TempDBClient) GetByKey(key string) (string, error) {
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

func (c *TempDBClient) SetEx(key string, seconds int, value string) (interface{}, error) {
	return c.sendCommand(fmt.Sprintf("SETEX %s %d %s", key, seconds, value))
}

func (c *TempDBClient) Delete(key string) (interface{}, error) {
	return c.sendCommand(fmt.Sprintf("DELETE %s", key))
}

func (c *TempDBClient) LPush(key, value string) (interface{}, error) {
	return c.sendCommand(fmt.Sprintf("LPUSH %s %s\r\n", key, value))
}

func (c *TempDBClient) SAdd(key, value string) (interface{}, error) {
	return c.sendCommand(fmt.Sprintf("SADD %s %s", key, value))
}

func (c *TempDBClient) Store(key string, value interface{}) (interface{}, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	log.Println("json value: ", string(jsonValue))
	return c.sendCommand(fmt.Sprintf("STORE %s %s", key, string(jsonValue)))
}

func (c *TempDBClient) StoreBatch(entries map[string]interface{}) (interface{}, error) {
	jsonValue, err := json.Marshal(entries)
	if err != nil {
		return "", err
	}
	return c.sendCommand(fmt.Sprintf("STOREBATCH %s", string(jsonValue)))
}

func (c *TempDBClient) GetFieldByKey(key, field string) (interface{}, error) {
	return c.sendCommand(fmt.Sprintf("GET_FIELD %s /%s", key, field))
}

func (c *TempDBClient) ViewData() (interface{}, error) {
	return c.sendCommand("VIEW_DATA")
}

// Command: SESSION_CREATE
func (c *TempDBClient) CreateSession(userID string) (interface{}, error) {
	command := fmt.Sprintf("SESSION_CREATE %s", userID)
	return c.sendCommand(command)
}

// Command: SESSION_GET
func (c *TempDBClient) GetSession(sessionID string) (interface{}, error) {
	command := fmt.Sprintf("SESSION_GET %s", sessionID)
	return c.sendCommand(command)
}

// Command: SESSION_SET
func (c *TempDBClient) SetSession(sessionID, key, value string) (interface{}, error) {
	command := fmt.Sprintf("SESSION_SET %s %s %s", sessionID, key, value)
	return c.sendCommand(command)
}

// Command: SESSION_DELETE
func (c *TempDBClient) DeleteSession(sessionID string) (interface{}, error) {
	command := fmt.Sprintf("SESSION_DELETE %s", sessionID)
	return c.sendCommand(command)
}
