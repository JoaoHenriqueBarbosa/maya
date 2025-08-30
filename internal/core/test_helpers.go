package core

import "time"

// MockEvent implements Event interface for testing
type MockEvent struct {
	EventType string
	EventData interface{}
	EventTime int64
}

func (m *MockEvent) Type() string {
	return m.EventType
}

func (m *MockEvent) Timestamp() int64 {
	if m.EventTime == 0 {
		return time.Now().UnixNano()
	}
	return m.EventTime
}