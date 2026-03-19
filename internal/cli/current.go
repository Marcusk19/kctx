package cli

import "fmt"

func runCurrent(short bool) {
	ctx := currentContext()
	if ctx == "" {
		fmt.Println("")
		return
	}
	fmt.Println(ctx)
}
