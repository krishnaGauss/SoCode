package models

import (
	"encoding/json"
	"time"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

type LogEntry struct {
	ID        string            `json:"id" db:"id"`
	Timestamp time.Time         `json:"timestamp" db:"timestamp"`
	Level     LogLevel          `json:"level" db:"level"`
	Message   string            `json:"message" db:"message"`
	Source    string            `json:"source" db:"source"`
	Service   string            `json:"service" db:"service"`
	Host      string            `json:"host" db:"host"`
	Tags      map[string]string `json:"tags" db:"tags"`
	Metadata  json.RawMessage   `json:"metadata" db:"metadata"`
}

type LogQuery struct {
	StartTime *time.Time        `json:"start_time,omitempty"`
	EndTime   *time.Time        `json:"end_time,omitempty"`
	Level     []LogLevel        `json:"level,omitempty"`
	Source    []string          `json:"source,omitempty"`
	Service   []string          `json:"service,omitempty"`
	Host      []string          `json:"host,omitempty"`
	Search    string            `json:"search,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	Limit     int               `json:"limit,omitempty"`
	Offset    int               `json:"offset,omitempty"`
}
