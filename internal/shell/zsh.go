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

oc() {
  if [ "$1" = "login" ]; then
    local _kctx_kubeconfig="$KUBECONFIG"
    export KUBECONFIG="${KCTX_ORIGINAL_KUBECONFIG:-${HOME}/.kube/config}"
    command oc "$@"
    local _kctx_rc=$?
    export KUBECONFIG="$_kctx_kubeconfig"
    [ $_kctx_rc -eq 0 ] && kctx _sync-context 2>/dev/null
    return $_kctx_rc
  fi
  command oc "$@"
}

trap '%s _cleanup-session 2>/dev/null' EXIT
`, "kctx", "kctx")
}
