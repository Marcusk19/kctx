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

func runCleanupSession() {
	sid := session.CurrentSessionID()
	if sid == "" {
		return
	}
	session.Remove(sid)
}
