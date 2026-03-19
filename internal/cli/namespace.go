package cli

import (
	"fmt"
	"os"

	"github.com/Marcusk19/kctx/internal/kubeconfig"
	"github.com/Marcusk19/kctx/internal/session"
)

func runNamespaceShow() {
	ctx := currentContext()
	if ctx == "" {
		fmt.Println("")
		return
	}

	// Check session overlay first for namespace override
	sid := session.CurrentSessionID()
	if sid != "" {
		cfg, err := kubeconfig.ReadConfig(session.ConfigPath(sid))
		if err == nil {
			for _, c := range cfg.Contexts {
				if c.Name == ctx && c.Context.Namespace != "" {
					fmt.Println(c.Context.Namespace)
					return
				}
			}
		}
	}

	// Fall back to source kubeconfig
	source := kubeconfig.SourceKubeconfig()
	ctxData := kubeconfig.GetContextInfo(source, ctx)
	if ctxData != nil && ctxData.Namespace != "" {
		fmt.Println(ctxData.Namespace)
	} else {
		fmt.Println("default")
	}
}

func runNamespaceSet(namespace string) {
	sid := requireSession()
	ctx := currentContext()
	if ctx == "" {
		fmt.Fprintln(os.Stderr, "error: no current context")
		os.Exit(1)
	}

	// Get the current context info from source
	source := kubeconfig.SourceKubeconfig()
	ctxData := kubeconfig.GetContextInfo(source, ctx)
	if ctxData == nil {
		// Check session overlay
		cfg, err := kubeconfig.ReadConfig(session.ConfigPath(sid))
		if err == nil {
			for _, c := range cfg.Contexts {
				if c.Name == ctx {
					ctxData = &c.Context
					break
				}
			}
		}
		if ctxData == nil {
			fmt.Fprintf(os.Stderr, "error: context %q not found\n", ctx)
			os.Exit(1)
		}
	}

	// Save previous namespace in metadata
	meta, _ := session.ReadMeta(sid)
	if meta == nil {
		meta = &session.Meta{}
	}
	// Get current namespace from session overlay first, then source
	currentNS := ""
	sessionCfg, err := kubeconfig.ReadConfig(session.ConfigPath(sid))
	if err == nil {
		for _, c := range sessionCfg.Contexts {
			if c.Name == ctx {
				currentNS = c.Context.Namespace
				break
			}
		}
	}
	if currentNS == "" {
		currentNS = ctxData.Namespace
	}
	if currentNS == "" {
		currentNS = "default"
	}
	meta.PreviousNS = currentNS
	meta.PreviousContext = ctx
	session.WriteMeta(sid, meta)

	if err := kubeconfig.WriteSessionConfigWithNamespace(session.Dir(sid), ctx, ctxData, namespace); err != nil {
		fmt.Fprintf(os.Stderr, "error setting namespace: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Namespace set to %q\n", namespace)
}

func runNamespacePrevious() {
	sid := requireSession()
	meta, err := session.ReadMeta(sid)
	if err != nil || meta.PreviousNS == "" {
		fmt.Fprintln(os.Stderr, "error: no previous namespace")
		os.Exit(1)
	}
	runNamespaceSet(meta.PreviousNS)
}
