package cli

import (
	"fmt"
	"os"

	"github.com/Marcusk19/kctx/internal/shell"
)

func runInit(shellName string) {
	switch shellName {
	case "zsh":
		fmt.Print(shell.ZshInit())
	case "bash":
		fmt.Print(shell.BashInit())
	case "fish":
		fmt.Print(shell.FishInit())
	default:
		fmt.Fprintf(os.Stderr, "unsupported shell: %s (supported: zsh, bash, fish)\n", shellName)
		os.Exit(1)
	}
}
