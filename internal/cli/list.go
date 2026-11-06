package cli

import (
	"fmt"
	"strings"

	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/store"
	"github.com/givikuna/ynternals/internal/types"
	"github.com/givikuna/ynternals/internal/validation"
)

func HandleList(inv types.Invocation) types.ExitOutcome {
	var secretsFile types.SecretsFilePath = types.SecretsFilePath(argAt(inv.Args, 0, string(types.DefaultSecretsFile)))

	if !store.FileExists(string(secretsFile)) {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: file not found at %s\n", secretsFile),
			ExitCode: 1,
		}
	}

	var rawRes fp.Result[[]byte] = store.ReadFileBytes(string(secretsFile))
	if rawRes.IsErr() {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: %s is not valid JSON\n", secretsFile),
			ExitCode: 1,
		}
	}
	rawBytes, _ := rawRes.Unwrap()

	var jsonRes fp.Result[any] = validation.ParseRawJSON(rawBytes)
	if jsonRes.IsErr() {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: %s is not valid JSON\n", secretsFile),
			ExitCode: 1,
		}
	}
	parsedJSON, _ := jsonRes.Unwrap()

	var objRes fp.Result[map[string]any] = validation.RequireJSONObject(parsedJSON)
	if objRes.IsErr() {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: %s must be a JSON object\n", secretsFile),
			ExitCode: 1,
		}
	}
	obj, _ := objRes.Unwrap()

	var keys []string = store.SortedKeys(obj)
	var output string = fp.Ternary(len(keys) == 0, "", strings.Join(keys, "\n")+"\n")

	return types.ExitOutcome{
		Stdout:   output,
		ExitCode: 0,
	}
}
