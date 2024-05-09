package domain

import "time"

type FileEvent struct {
	Path      string
	Timestamp time.Time
	Type      EventType
}

type EventType string

const (
	Create  EventType = "CREATE"
	Write   EventType = "WRITE"
	Remove  EventType = "REMOVE"
	Rename  EventType = "RENAME"
	Close   EventType = "CLOSE"
	Unknown EventType = "UNKNOWN"
)
