package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/Marcusk19/kctx/internal/kubeconfig"
	"github.com/Marcusk19/kctx/internal/session"
)

func runInitSession() {
	sid := session.CurrentSessionID()
	if sid == "" {
		os.Exit(1)
	}

	if session.Exists(sid) {
		return
	}

	// Get current context from source kubeconfig to inherit
	source := kubeconfig.SourceKubeconfig()
	_, currentCtx, _ := kubeconfig.ListContexts(source)

	if err := kubeconfig.WriteSessionConfig(session.Dir(sid), currentCtx); err != nil {
		fmt.Fprintf(os.Stderr, "error creating session: %v\n", err)
		os.Exit(1)
	}

	meta := &session.Meta{
		CreatedAt: time.Now(),
	}
	session.WriteMeta(sid, meta)
}

func runSyncContext() {
	sid := requireSession()

	// Read current context from source kubeconfig
	source := kubeconfig.SourceKubeconfig()
	_, sourceCtx, err := kubeconfig.ListContexts(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading source kubeconfig: %v\n", err)
		os.Exit(1)
	}
	if sourceCtx == "" {
		fmt.Fprintln(os.Stderr, "error: no current-context in source kubeconfig")
		os.Exit(1)
	}

	// Save previous context in metadata
	prev := currentContext()
	meta, _ := session.ReadMeta(sid)
	if meta == nil {
		meta = &session.Meta{CreatedAt: time.Now()}
	}
	meta.PreviousContext = prev
	session.WriteMeta(sid, meta)

	// Rewrite session overlay to point to the source's current-context
	if err := kubeconfig.WriteSessionConfig(session.Dir(sid), sourceCtx); err != nil {
		fmt.Fprintf(os.Stderr, "error syncing context: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Synced to context %q\n", sourceCtx)
}

func runCleanupSession() {
	sid := session.CurrentSessionID()
	if sid == "" {
		return
	}
	session.Remove(sid)
}
