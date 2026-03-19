package cli

import (
	"fmt"
	"os"
)

var version = "dev"

func Run() {
	args := os.Args[1:]

	if len(args) == 0 {
		runList()
		return
	}

	switch args[0] {
	case "init":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: kctx init <zsh|bash|fish>")
			os.Exit(1)
		}
		runInit(args[1])
	case "current":
		short := hasFlag(args, "--short")
		runCurrent(short)
	case "ns":
		if len(args) < 2 {
			runNamespaceShow()
		} else if args[1] == "-" {
			runNamespacePrevious()
		} else {
			runNamespaceSet(args[1])
		}
	case "cleanup":
		runCleanup()
	case "version":
		fmt.Printf("kctx %s\n", version)
	case "_init-session":
		runInitSession()
	case "_cleanup-session":
		runCleanupSession()
	case "-":
		runSwitchPrevious()
	case "--help", "-h", "help":
		printUsage()
	default:
		runSwitch(args[0])
	}
}

func printUsage() {
	fmt.Print(`kctx - per-shell Kubernetes context isolation

Usage:
  kctx                     List all contexts (current highlighted)
  kctx <name>              Switch to context (shell-local)
  kctx -                   Switch to previous context
  kctx current             Show current context
  kctx ns                  Show current namespace
  kctx ns <namespace>      Set namespace (shell-local)
  kctx ns -                Switch to previous namespace
  kctx init <shell>        Print shell integration code (zsh|bash|fish)
  kctx cleanup             Remove stale session configs
  kctx version             Print version

Environment Variables:
  KCTX_SESSION              Session ID (set by shell integration)
  KCTX_SOURCE_KUBECONFIG    Override source kubeconfig
  KCTX_DIR                  Override data directory (default: ~/.kctx)
`)
}

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}
