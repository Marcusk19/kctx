package cli

import (
	"fmt"
	"os"

	"github.com/Marcusk19/kctx/internal/kubeconfig"
	"github.com/Marcusk19/kctx/internal/session"
)

func runSwitch(name string) {
	sid := requireSession()

	// Validate context exists in source
	source := kubeconfig.SourceKubeconfig()
	contexts, _, err := kubeconfig.ListContexts(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading kubeconfig: %v\n", err)
		os.Exit(1)
	}
	found := false
	for _, c := range contexts {
		if c == name {
			found = true
			break
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "error: context %q not found\n", name)
		os.Exit(1)
	}

	// Save previous context in metadata
	prev := currentContext()
	meta, _ := session.ReadMeta(sid)
	if meta == nil {
		meta = &session.Meta{}
	}
	meta.PreviousContext = prev
	session.WriteMeta(sid, meta)

	// Write session config
	if err := kubeconfig.WriteSessionConfig(session.Dir(sid), name); err != nil {
		fmt.Fprintf(os.Stderr, "error switching context: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Switched to context %q\n", name)
}

func runSwitchPrevious() {
	sid := requireSession()
	meta, err := session.ReadMeta(sid)
	if err != nil || meta.PreviousContext == "" {
		fmt.Fprintln(os.Stderr, "error: no previous context")
		os.Exit(1)
	}
	runSwitch(meta.PreviousContext)
}

func requireSession() string {
	sid := session.CurrentSessionID()
	if sid == "" {
		fmt.Fprintln(os.Stderr, "error: KCTX_SESSION not set. Run: eval \"$(kctx init zsh)\"")
		os.Exit(1)
	}
	return sid
}
