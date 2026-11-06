package cli

import (
	"fmt"

	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/store"
	"github.com/givikuna/ynternals/internal/types"
)

func HandleRemove(inv types.Invocation) types.ExitOutcome {
	var name string = argAt(inv.Args, 0, "")
	var secFile types.SecretsFilePath = types.SecretsFilePath(argAt(inv.Args, 1, string(types.DefaultSecretsFile)))

	if name == "" {
		return types.ExitOutcome{
			Stdout:   "usage: ynternals remove <secret_name> [secrets.json]\n",
			ExitCode: 1,
		}
	}

	if !store.FileExists(string(secFile)) {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: file not found at %s\n", secFile),
			ExitCode: 1,
		}
	}

	var objRes fp.Result[map[string]any] = loadAndValidateObjectAborting(secFile)
	if objRes.IsErr() {
		_, err := objRes.Unwrap()
		return types.ExitOutcome{
			Stdout:   err.Error(),
			ExitCode: 1,
		}
	}
	var existingObj map[string]any
	existingObj, _ = objRes.Unwrap()

	var updatedObj map[string]any = store.DeleteSecret(existingObj, types.SecretName(name))
	var serializedRes fp.Result[[]byte] = store.SerializeJSONObject(updatedObj)

	if serializedRes.IsErr() {
		_, err := serializedRes.Unwrap()
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: %v\n", err),
			ExitCode: 1,
		}
	}
	var serializedBytes []byte
	serializedBytes, _ = serializedRes.Unwrap()

	var err error
	if err = store.WriteAtomicTemp(string(secFile), serializedBytes); err != nil {
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: %v\n", err),
			ExitCode: 1,
		}
	}

	return types.ExitOutcome{
		Stdout:   fmt.Sprintf("removed '%s' from %s\n", name, secFile),
		ExitCode: 0,
	}
}
