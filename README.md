# kctx

Per-shell Kubernetes context isolation. Switch contexts and namespaces in one terminal without affecting others.

## The Problem

With standard `kubectl`, switching context (`kubectl config use-context`) modifies `~/.kube/config` globally. Every open terminal immediately sees the change, which can lead to running commands against the wrong cluster.

## How kctx Solves It

kctx gives each shell session its own lightweight kubeconfig overlay. When you switch context with kctx, it writes only to your session's overlay file (`~/.kctx/sessions/<session-id>/config`), which is prepended to `KUBECONFIG`. Kubernetes merges configs left-to-right, so your session's `current-context` takes priority without touching the real kubeconfig.

Session IDs are derived automatically from tmux pane, SSH TTY, or shell PID.

## Installation

### From source

Requires Go 1.25+.

```sh
git clone https://github.com/Marcusk19/kctx.git
cd kctx
make install    # installs to $GOPATH/bin
```

### Shell integration

Add one of the following to your shell's rc file:

**zsh** (`~/.zshrc`):
```sh
eval "$(kctx init zsh)"
```

**bash** (`~/.bashrc`):
```sh
eval "$(kctx init bash)"
```

**fish** (`~/.config/fish/config.fish`):
```fish
kctx init fish | source
```

Shell integration sets up the session, configures `KUBECONFIG`, wraps `oc login` to write to the original kubeconfig (then sync the session), and registers a cleanup trap on shell exit.

## Usage

```
kctx                     List all contexts (current highlighted)
kctx <name>              Switch to context (shell-local)
kctx -                   Switch to previous context
kctx current             Show current context
kctx ns                  Show current namespace
kctx ns <namespace>      Set namespace (shell-local)
kctx ns -                Switch to previous namespace
kctx cleanup             Remove stale session configs
kctx version             Print version
```

### Examples

```sh
# List available contexts — current one is highlighted green
$ kctx
* dev-cluster
  staging-cluster
  prod-cluster

# Switch context in this shell only
$ kctx staging-cluster
Switched to context "staging-cluster"

# Switch back to previous context
$ kctx -
Switched to context "dev-cluster"

# Set namespace for this shell
$ kctx ns my-app
Namespace set to "my-app"

# Switch back to previous namespace
$ kctx ns -
Namespace set to "default"

# Clean up sessions from closed terminals
$ kctx cleanup
Removed session: shell-12345
Cleaned 1 stale session(s).
```

## Environment Variables

| Variable | Description |
|---|---|
| `KCTX_SESSION` | Session ID (set automatically by shell integration) |
| `KCTX_SOURCE_KUBECONFIG` | Override the source kubeconfig path |
| `KCTX_DIR` | Override the data directory (default: `~/.kctx`) |

## Building & Development

```sh
make build      # Build binary to ./kctx
make test       # Run tests
make install    # Install to $GOPATH/bin
make clean      # Remove built binary
```

Version is embedded at build time via ldflags from `git describe --tags`.

## License

Apache License 2.0
