package cli

import (
	"fmt"

	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/store"
	"github.com/givikuna/ynternals/internal/types"
	"github.com/givikuna/ynternals/internal/validation"
)

func HandleValidate(inv types.Invocation) types.ExitOutcome {
	var secFile types.SecretsFilePath = types.SecretsFilePath(argAt(inv.Args, 0, string(types.DefaultSecretsFile)))

	if !store.FileExists(string(secFile)) {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: file not found at %s\n", secFile),
			ExitCode: 1,
		}
	}

	var rawRes fp.Result[[]byte] = store.ReadFileBytes(string(secFile))
	if rawRes.IsErr() {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: %s is not valid JSON\n", secFile),
			ExitCode: 1,
		}
	}
	var rawBytes []byte
	rawBytes, _ = rawRes.Unwrap()

	var jsonRes fp.Result[any] = validation.ParseRawJSON(rawBytes)
	if jsonRes.IsErr() {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: %s is not valid JSON\n", secFile),
			ExitCode: 1,
		}
	}
	parsedJSON, _ := jsonRes.Unwrap()

	var objRes fp.Result[map[string]any] = validation.RequireJSONObject(parsedJSON)
	if objRes.IsErr() {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: %s must be a JSON object\n", secFile),
			ExitCode: 1,
		}
	}
	var obj map[string]any
	obj, _ = objRes.Unwrap()

	var strictRes fp.Result[types.SecretMap] = validation.RequireStringValues(obj)
	if strictRes.IsErr() {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: all values in %s must be strings\n", secFile),
			ExitCode: 1,
		}
	}

	return types.ExitOutcome{
		Stdout:   fmt.Sprintf("%s is valid.\n", secFile),
		ExitCode: 0,
	}
}
