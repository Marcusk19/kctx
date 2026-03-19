package session

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// DefaultTTL is the default session time-to-live.
const DefaultTTL = 7 * 24 * time.Hour

// CleanupStale removes sessions that are older than the TTL or whose
// owning process no longer exists.
func CleanupStale(ttl time.Duration) ([]string, error) {
	sessionsDir := filepath.Join(BaseDir(), "sessions")
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var cleaned []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		shouldClean := false

		// Check if the owning process is still alive
		if strings.HasPrefix(name, "shell-") {
			pidStr := strings.TrimPrefix(name, "shell-")
			if pid, err := strconv.Atoi(pidStr); err == nil {
				if !processExists(pid) {
					shouldClean = true
				}
			}
		}

		// Check TTL via meta.json
		if !shouldClean {
			meta, err := ReadMeta(name)
			if err != nil {
				// No meta.json — check dir mtime
				info, err := e.Info()
				if err == nil && time.Since(info.ModTime()) > ttl {
					shouldClean = true
				}
			} else if time.Since(meta.LastUsed) > ttl {
				shouldClean = true
			}
		}

		if shouldClean {
			if err := Remove(name); err == nil {
				cleaned = append(cleaned, name)
			}
		}
	}
	return cleaned, nil
}

func processExists(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 checks if process exists without actually sending a signal
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}
