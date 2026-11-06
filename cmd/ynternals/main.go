package main

import (
	"fmt"
	"os"

	"github.com/givikuna/ynternals/internal/cli"
)

var version = "dev"

func main() {
	inv := cli.BuildInvocation(os.Args[1:], version)
	outcome := cli.Dispatch(inv)

	if outcome.Stdout != "" {
		_, _ = fmt.Fprint(os.Stdout, outcome.Stdout)
	}
	if outcome.Stderr != "" {
		_, _ = fmt.Fprint(os.Stderr, outcome.Stderr)
	}
	os.Exit(outcome.ExitCode)
}
