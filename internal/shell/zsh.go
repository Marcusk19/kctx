package shell

import "fmt"

func ZshInit() string {
	return fmt.Sprintf(`# kctx shell integration for zsh
export KCTX_ORIGINAL_KUBECONFIG="${KUBECONFIG}"

if [ -n "$TMUX_PANE" ]; then
  export KCTX_SESSION="tmux-${TMUX_PANE}"
elif [ -n "$SSH_TTY" ]; then
  export KCTX_SESSION="ssh-$(echo "$SSH_TTY" | shasum | cut -c1-8)"
else
  export KCTX_SESSION="shell-$$"
fi

%s _init-session 2>/dev/null

export KUBECONFIG="${HOME}/.kctx/sessions/${KCTX_SESSION}/config:${KCTX_ORIGINAL_KUBECONFIG:-${HOME}/.kube/config}"

trap '%s _cleanup-session 2>/dev/null' EXIT
`, "kctx", "kctx")
}
