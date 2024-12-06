package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Addr       string
	Collection string
}

type TempDBClient struct {
	conn      net.Conn
	addr      string
	collection string
	mu        sync.Mutex
}

type clientPool struct {
	clients chan *TempDBClient
	size    int
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
	return &TempDBClient{conn: conn, addr: config.Addr, collection: config.Collection}, nil
}

func (c *TempDBClient) Close() {
	if pool.size > len(pool.clients) {
		pool.clients <- c
	} else {
		c.conn.Close()
	}
}

func (c *TempDBClient) Ping() (string, error) {
	pong, err := c.sendCommand("PING")
	return pong, err
}

func Client(addr, collection string) (*TempDBClient, error) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	client := &TempDBClient{conn: conn, addr: addr, collection: collection}

	return client, nil
}

func (c *TempDBClient) sendCommand(command string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fullCommand := fmt.Sprintf("%s %s", c.collection, command)
	_, err := fmt.Fprintf(c.conn, fullCommand+"\r\n")
	if err != nil {
		return "", err
	}

	response, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

func (c *TempDBClient) Set(key, value string) (string, error) {
	return c.sendCommand(fmt.Sprintf("SET %s %s", key, value))
}

func (c *TempDBClient) GetByKey(key string) (string, error) {
	return c.sendCommand(fmt.Sprintf("GET_KEY %s", key))
}

func (c *TempDBClient) SetEx(key string, seconds int, value string) (string, error) {
	return c.sendCommand(fmt.Sprintf("SETEX %s %d %s", key, seconds, value))
}

func (c *TempDBClient) Delete(key string) (string, error) {
	return c.sendCommand(fmt.Sprintf("DELETE %s", key))
}

func (c *TempDBClient) LPush(key, value string) (string, error) {
	return c.sendCommand(fmt.Sprintf("LPUSH %s %s\r\n", key, value))
}

func (c *TempDBClient) SAdd(key, value string) (string, error) {
	return c.sendCommand(fmt.Sprintf("SADD %s %s", key, value))
}

func (c *TempDBClient) Store(key string, value interface{}) (string, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return c.sendCommand(fmt.Sprintf("STORE %s %s", key, string(jsonValue)))
}

func (c *TempDBClient) GetFieldByKey(key, field string) (string, error) {
	return c.sendCommand(fmt.Sprintf("GET_FIELD %s /%s", key, field))
}

func (c *TempDBClient) ViewData() (string, error) {
	return c.sendCommand("VIEW_DATA")
}

func (c *TempDBClient) GetDB() (string, error) {
	return c.sendCommand("GET_DB")
}
