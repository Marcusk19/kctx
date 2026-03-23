# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```sh
make build        # Build binary to ./kctx
make install      # Install to $GOPATH/bin
make test         # Run all tests: go test ./...
make clean        # Remove built binary
```

Version is injected via `-ldflags` from `git describe --tags`.

## Architecture

kctx provides **per-shell Kubernetes context isolation**. Each shell session gets its own kubeconfig overlay (stored in `~/.kctx/sessions/<session-id>/config`) that is prepended to `KUBECONFIG`, so context switches in one terminal don't affect others.

### How it works

1. **Shell integration** (`kctx init zsh|bash|fish`) outputs shell code that:
   - Sets `KCTX_SESSION` to a unique session ID (based on tmux pane, SSH TTY, or shell PID)
   - Saves original `KUBECONFIG` as `KCTX_ORIGINAL_KUBECONFIG`
   - Prepends the session overlay path to `KUBECONFIG`
   - Wraps `oc login` to write to the original kubeconfig then sync the session
   - Registers a cleanup trap on shell exit

2. **Session overlay** is a minimal kubeconfig YAML containing only `current-context` (and optionally a context entry for namespace overrides). Kubernetes merges this with the real kubeconfig via the colon-separated `KUBECONFIG` path.

3. **Metadata** (`meta.json`) tracks previous context/namespace for the `-` (switch back) feature.

### Package structure

- `cmd/kctx/main.go` — Entry point, calls `cli.Run()`
- `internal/cli/` — Command dispatch (`root.go`) and handlers. No flag library; uses manual `os.Args` parsing. Internal commands prefixed with `_` are called by shell integration, not users.
- `internal/kubeconfig/` — Read/write kubeconfig YAML files. `types.go` defines minimal kubeconfig structs.
- `internal/session/` — Session directory management and metadata (JSON).
- `internal/shell/` — Shell init script generators (one file per shell).

### Key environment variables

- `KCTX_SESSION` — Session ID (set by shell integration)
- `KCTX_SOURCE_KUBECONFIG` / `KCTX_ORIGINAL_KUBECONFIG` — Override source kubeconfig path
- `KCTX_DIR` — Override data directory (default: `~/.kctx`)
