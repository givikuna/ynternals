package cli

import (
	"fmt"
	"os"

	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/store"
	"github.com/givikuna/ynternals/internal/types"
	"github.com/givikuna/ynternals/internal/validation"
)

var commandRoutes = map[types.CommandName]types.CommandHandler{
	types.CmdVersionShort: HandleVersion,
	types.CmdVersionLong:  HandleVersion,
	types.CmdAdd:          HandleAdd,
	types.CmdRemove:       HandleRemove,
	types.CmdValidate:     HandleValidate,
	types.CmdList:         HandleList,
	types.CmdPeek:         HandlePeek,
}

func BuildInvocation(rawArgs []string, version string) types.Invocation {
	var envKeyFile string = os.Getenv(types.KeyFileEnvVar)

	var keyFile types.KeyFilePath = fp.Ternary(
		envKeyFile != "",
		types.KeyFilePath(envKeyFile),
		types.DefaultKeyFile,
	)

	var cmdName types.CommandName = fp.Ternary(len(rawArgs) > 0, types.CommandName(rawArgs[0]), types.CommandName(""))
	var restArgs []string = fp.Ternary(len(rawArgs) > 1, rawArgs[1:], []string(nil))

	return types.Invocation{
		Command: cmdName,
		Args:    restArgs,
		Env: types.EnvConfig{
			KeyFile:    keyFile,
			RuntimeDir: types.DefaultRuntimeDir,
			EUID:       os.Geteuid(),
			Version:    version,
		},
	}
}

func Dispatch(inv types.Invocation) types.ExitOutcome {
	var handler types.CommandHandler
	var found bool
	handler, found = commandRoutes[inv.Command]

	return fp.Ternary(found, handler, HandleUsage)(inv)
}

func argAt(args []string, index int, fallback string) string {
	return fp.Ternary(index >= 0 && index < len(args) && args[index] != "", args[index], fallback)
}

func loadAndValidateObjectAborting(secretsFile types.SecretsFilePath) fp.Result[map[string]any] {
	var rawRes fp.Result[[]byte] = store.ReadFileBytes(string(secretsFile))
	if rawRes.IsErr() {
		return fp.Err[map[string]any](fmt.Errorf("ERROR: %s is not valid JSON. Aborting.\n", secretsFile))
	}
	rawBytes, _ := rawRes.Unwrap()

	var jsonRes fp.Result[any] = validation.ParseRawJSON(rawBytes)
	if jsonRes.IsErr() {
		return fp.Err[map[string]any](fmt.Errorf("ERROR: %s is not valid JSON. Aborting.\n", secretsFile))
	}
	parsed, _ := jsonRes.Unwrap()

	var objRes fp.Result[map[string]any] = validation.RequireJSONObject(parsed)
	if objRes.IsErr() {
		return fp.Err[map[string]any](fmt.Errorf("ERROR: %s must be a JSON object. Aborting.\n", secretsFile))
	}
	return objRes
}
