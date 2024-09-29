package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type LogEntry struct {
	Timestamp     string            `json:"timestamp"`
	Message       string            `json:"message"`
	Level         LogLevel          `json:"level"`
	ApplicationID string            `json:"application_id"`
	Tags          map[string]string `json:"tags"`
}

type LoggerClient interface {
	Log(level LogLevel, message string, tags map[string]string) error
	GetLogs() (map[LogLevel][]LogEntry, error)
}

type HttpLoggerClient struct {
	ApplicationID string
	ServerURL     string
}

func NewHttpLoggerClient(applicationID, serverURL string) *HttpLoggerClient {
	return &HttpLoggerClient{
		ApplicationID: applicationID,
		ServerURL:     serverURL,
	}
}

func (c *HttpLoggerClient) Log(level LogLevel, message string, tags map[string]string) error {
	entry := LogEntry{
		Timestamp:     time.Now().Format(time.RFC3339),
		Message:       message,
		Level:         level,
		ApplicationID: c.ApplicationID,
		Tags:          tags,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.ServerURL+"/log", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send log: status code %d", resp.StatusCode)
	}

	return nil
}

func (c *HttpLoggerClient) GetLogs() (map[LogLevel][]LogEntry, error) {
	resp, err := http.Get(fmt.Sprintf("%s/getLogs?appId=%s", c.ServerURL, c.ApplicationID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get logs: status code %d", resp.StatusCode)
	}

	var logs map[LogLevel][]LogEntry
	err = json.NewDecoder(resp.Body).Decode(&logs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
