package cli

import (
	"fmt"

	"github.com/givikuna/ynternals/internal/types"
)

const usageText = `usage:
  ynternals add
  ynternals remove <secret_name> [secrets.json path]
  ynternals validate [secrets.json path]
  ynternals list [secrets.json path]
  ynternals peek <secret_name>
  ynternals -v, --version
`

func HandleVersion(inv types.Invocation) types.ExitOutcome {
	return types.ExitOutcome{
		Stdout:   fmt.Sprintf("ynternals version %s\n", inv.Env.Version),
		ExitCode: 0,
	}
}

func HandleUsage(_ types.Invocation) types.ExitOutcome {
	return types.ExitOutcome{
		Stdout:   usageText,
		ExitCode: 1,
	}
}
