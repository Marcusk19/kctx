package session

import (
	"os"
	"path/filepath"
)

// Dir returns the session directory for the given session ID.
func Dir(sessionID string) string {
	return filepath.Join(BaseDir(), "sessions", sessionID)
}

// BaseDir returns the kctx data directory.
func BaseDir() string {
	if v := os.Getenv("KCTX_DIR"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kctx")
}

// ConfigPath returns the path to the session's kubeconfig overlay.
func ConfigPath(sessionID string) string {
	return filepath.Join(Dir(sessionID), "config")
}

// CurrentSessionID returns the current session ID from the environment.
func CurrentSessionID() string {
	return os.Getenv("KCTX_SESSION")
}

// Exists returns true if the session directory exists.
func Exists(sessionID string) bool {
	_, err := os.Stat(Dir(sessionID))
	return err == nil
}

// Remove removes the session directory.
func Remove(sessionID string) error {
	return os.RemoveAll(Dir(sessionID))
}
