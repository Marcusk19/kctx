package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Meta holds session metadata.
type Meta struct {
	CreatedAt       time.Time `json:"created_at"`
	LastUsed        time.Time `json:"last_used"`
	PreviousContext string    `json:"previous_context,omitempty"`
	PreviousNS      string    `json:"previous_ns,omitempty"`
}

// ReadMeta reads the session metadata.
func ReadMeta(sessionID string) (*Meta, error) {
	path := filepath.Join(Dir(sessionID), "meta.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Meta
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// WriteMeta writes the session metadata.
func WriteMeta(sessionID string, m *Meta) error {
	if err := os.MkdirAll(Dir(sessionID), 0700); err != nil {
		return err
	}
	m.LastUsed = time.Now()
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(Dir(sessionID), "meta.json"), data, 0600)
}
