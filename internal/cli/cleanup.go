package cli

import (
	"fmt"
	"os"

	"github.com/Marcusk19/kctx/internal/session"
)

func runCleanup() {
	cleaned, err := session.CleanupStale(session.DefaultTTL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error during cleanup: %v\n", err)
		os.Exit(1)
	}
	if len(cleaned) == 0 {
		fmt.Println("No stale sessions found.")
		return
	}
	for _, s := range cleaned {
		fmt.Printf("Removed session: %s\n", s)
	}
	fmt.Printf("Cleaned %d stale session(s).\n", len(cleaned))
}
