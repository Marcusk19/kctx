package shell

import "fmt"

func FishInit() string {
	return fmt.Sprintf(`# kctx shell integration for fish
set -gx KCTX_ORIGINAL_KUBECONFIG "$KUBECONFIG"

if set -q TMUX_PANE
  set -gx KCTX_SESSION "tmux-$TMUX_PANE"
else if set -q SSH_TTY
  set -gx KCTX_SESSION "ssh-"(echo "$SSH_TTY" | shasum | cut -c1-8)
else
  set -gx KCTX_SESSION "shell-$fish_pid"
end

%s _init-session 2>/dev/null

if test -n "$KCTX_ORIGINAL_KUBECONFIG"
  set -gx KUBECONFIG "$HOME/.kctx/sessions/$KCTX_SESSION/config:$KCTX_ORIGINAL_KUBECONFIG"
else
  set -gx KUBECONFIG "$HOME/.kctx/sessions/$KCTX_SESSION/config:$HOME/.kube/config"
end

function oc
  if test (count $argv) -gt 0; and test "$argv[1]" = "login"
    set -l _kctx_kubeconfig $KUBECONFIG
    set -gx KUBECONFIG (test -n "$KCTX_ORIGINAL_KUBECONFIG"; and echo "$KCTX_ORIGINAL_KUBECONFIG"; or echo "$HOME/.kube/config")
    command oc $argv
    set -l _kctx_rc $status
    set -gx KUBECONFIG "$_kctx_kubeconfig"
    test $_kctx_rc -eq 0; and kctx _sync-context 2>/dev/null
    return $_kctx_rc
  end
  command oc $argv
end

function __kctx_cleanup --on-event fish_exit
  %s _cleanup-session 2>/dev/null
end
`, "kctx", "kctx")
}
