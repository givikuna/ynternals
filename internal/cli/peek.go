package cli

import (
	"fmt"
	"path/filepath"

	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/store"
	"github.com/givikuna/ynternals/internal/types"
)

func HandlePeek(inv types.Invocation) types.ExitOutcome {
	var name string = argAt(inv.Args, 0, "")

	if name == "" {
		return types.ExitOutcome{
			Stdout:   "usage: ynternals peek <secret_name>\n",
			ExitCode: 1,
		}
	}

	if inv.Env.EUID != 0 {
		return types.ExitOutcome{
			Stdout:   "ERROR: peeking secrets requires root privileges (use sudo)\n",
			ExitCode: 1,
		}
	}

	var secretPath string = filepath.Join(string(inv.Env.RuntimeDir), name)
	if !store.FileExists(secretPath) {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: secret '%s' not found in %s\n", name, secretPath),
			ExitCode: 1,
		}
	}

	var contentRes fp.Result[[]byte] = store.ReadFileBytes(secretPath)
	if contentRes.IsErr() {
		_, err := contentRes.Unwrap()

		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: %v\n", err),
			ExitCode: 1,
		}
	}
	content, _ := contentRes.Unwrap()

	return types.ExitOutcome{
		Stdout:   string(content) + "\n",
		ExitCode: 0,
	}
}
