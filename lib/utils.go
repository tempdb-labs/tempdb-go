// Package lib provides utility functions and types for handling and formatting responses.
package lib

import (
	"encoding/json"
	"fmt"
)

// Response represents a standard response structure with a status, message, and data.
type Response struct {
	Status  string          `json:"status"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data"`
}

// ResponseData represents the data structure within a response, including a type and raw data.
type ResponseData struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// FormattedResponse represents a formatted response with additional metadata such as type, value, creation time, and expiry time.
type FormattedResponse struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	CreatedAt *int64      `json:"created_at,omitempty"`
	Expiry    *int64      `json:"expiry,omitempty"`
}

// formatResponse formats the given data as a pretty-printed JSON string.
// It returns the formatted JSON string or an error if the formatting fails.
func formatResponse(data interface{}) (string, error) {
	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error formatting response: %w", err)
	}
	return string(formatted), nil
}
