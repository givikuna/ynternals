package main

import (
	"fmt"
	"os"

	"github.com/givikuna/ynternals/internal/decrypt"
	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/types"
)

func main() {
	secretsArg := fp.Ternary(len(os.Args) > 1, os.Args[1], "")
	keyArg := fp.Ternary(len(os.Args) > 2, os.Args[2], "")

	outcome := decrypt.Run(
		types.SecretsFilePath(secretsArg),
		types.KeyFilePath(keyArg),
		types.DefaultRuntimeDir,
	)

	if outcome.Stdout != "" {
		_, _ = fmt.Fprint(os.Stdout, outcome.Stdout)
	}
	if outcome.Stderr != "" {
		_, _ = fmt.Fprint(os.Stderr, outcome.Stderr)
	}
	os.Exit(outcome.ExitCode)
}
