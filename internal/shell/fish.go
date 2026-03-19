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

function __kctx_cleanup --on-event fish_exit
  %s _cleanup-session 2>/dev/null
end
`, "kctx", "kctx")
}
