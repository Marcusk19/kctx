package cli

import (
	"fmt"
	"os"

	"github.com/Marcusk19/kctx/internal/kubeconfig"
	"github.com/Marcusk19/kctx/internal/session"
)

func runList() {
	source := kubeconfig.SourceKubeconfig()
	contexts, _, err := kubeconfig.ListContexts(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing contexts: %v\n", err)
		os.Exit(1)
	}

	current := currentContext()

	isTTY := isTerminal(os.Stdout)
	for _, ctx := range contexts {
		if ctx == current {
			if isTTY {
				fmt.Printf("\033[1;32m* %s\033[0m\n", ctx)
			} else {
				fmt.Printf("* %s\n", ctx)
			}
		} else {
			fmt.Printf("  %s\n", ctx)
		}
	}
}

func currentContext() string {
	// First check the session overlay
	sid := session.CurrentSessionID()
	if sid != "" {
		cfg, err := kubeconfig.ReadConfig(session.ConfigPath(sid))
		if err == nil && cfg.CurrentContext != "" {
			return cfg.CurrentContext
		}
	}
	// Fall back to source kubeconfig
	source := kubeconfig.SourceKubeconfig()
	_, current, _ := kubeconfig.ListContexts(source)
	return current
}
