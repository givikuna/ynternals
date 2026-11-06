package decrypt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/givikuna/ynternals/internal/crypto"
	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/store"
	"github.com/givikuna/ynternals/internal/types"
	"github.com/givikuna/ynternals/internal/validation"
)

func Run(secretsFile types.SecretsFilePath, keyFile types.KeyFilePath, runtimeDir types.RuntimeDirPath) types.ExitOutcome {
	if !store.FileExists(string(secretsFile)) || !store.FileExists(string(keyFile)) {
		return types.ExitOutcome{
			Stdout:   "secrets or key file missing.\n",
			ExitCode: 0,
		}
	}

	var rawRes fp.Result[[]byte] = store.ReadFileBytes(string(secretsFile))
	if rawRes.IsErr() {
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: %s is not valid JSON\n", secretsFile),
			ExitCode: 1,
		}
	}
	rawBytes, _ := rawRes.Unwrap()

	var jsonRes fp.Result[any] = validation.ParseRawJSON(rawBytes)
	if jsonRes.IsErr() {
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: %s is not valid JSON\n", secretsFile),
			ExitCode: 1,
		}
	}
	parsedJSON, _ := jsonRes.Unwrap()

	var objRes fp.Result[map[string]any] = validation.RequireJSONObject(parsedJSON)
	if objRes.IsErr() {
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: %s must be a JSON object\n", secretsFile),
			ExitCode: 1,
		}
	}
	obj, _ := objRes.Unwrap()

	var secretsRes fp.Result[types.SecretMap] = validation.RequireStringValues(obj)
	if secretsRes.IsErr() {
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: all values in %s must be strings\n", secretsFile),
			ExitCode: 1,
		}
	}
	secretMap, _ := secretsRes.Unwrap()

	var passRes fp.Result[types.Passphrase] = store.
		ReadPassphraseFile(keyFile)
	if passRes.IsErr() {
		_, err := passRes.Unwrap()
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("%v\n", err),
			ExitCode: 1,
		}
	}
	passphrase, _ := passRes.Unwrap()

	var entries []types.SecretEntry = fp.Filter(store.SortedEntries(secretMap), func(e types.SecretEntry) bool {
		return e.Name != ""
	})

	var err error = fp.ForEachErr(entries, func(entry types.SecretEntry) error {
		if !validation.IsValidSecretName(string(entry.Name)) {
			return fmt.Errorf("invalid secret name: %s", entry.Name)
		}

		var plainRes fp.Result[types.Plaintext] = crypto.DecryptAES256CBC(entry.Value, passphrase)
		if plainRes.IsErr() {
			_, decErr := plainRes.Unwrap()
			return decErr
		}
		plaintext, _ := plainRes.Unwrap()

		var targetPath string = filepath.Join(string(runtimeDir), string(entry.Name))
		return os.WriteFile(targetPath, plaintext, 0o666)
	})

	if err != nil {
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("%v\n", err),
			ExitCode: 1,
		}
	}

	return types.ExitOutcome{ExitCode: 0}
}
