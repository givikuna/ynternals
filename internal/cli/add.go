package cli

import (
	"fmt"
	"os"

	"github.com/givikuna/ynternals/internal/crypto"
	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/store"
	"github.com/givikuna/ynternals/internal/types"
	"github.com/givikuna/ynternals/internal/validation"
)

func HandleAdd(inv types.Invocation) types.ExitOutcome {
	var name string = argAt(inv.Args, 0, "")
	var secretsFile types.SecretsFilePath = types.SecretsFilePath(argAt(inv.Args, 1, string(types.DefaultSecretsFile)))

	if name == "" {
		return types.ExitOutcome{
			Stdout:   "usage: ynternals add <secret_name> [secrets.json path]\n",
			ExitCode: 1,
		}
	}

	var nameRes fp.Result[types.SecretName] = validation.ParseSecretName(name)
	if nameRes.IsErr() {
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: invalid secret name: %s\n", name),
			ExitCode: 1,
		}
	}

	secretName, _ := nameRes.Unwrap()

	if !store.FileExists(string(inv.Env.KeyFile)) {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: key not found at %s\n", inv.Env.KeyFile),
			ExitCode: 1,
		}
	}

	if err := store.EnsureSecretsFileExists(secretsFile); err != nil {
		return types.ExitOutcome{
			Stdout:   fmt.Sprintf("ERROR: %v\n", err),
			ExitCode: 1,
		}
	}

	var objRes fp.Result[map[string]any] = loadAndValidateObjectAborting(secretsFile)
	if objRes.IsErr() {
		_, err := objRes.Unwrap()
		return types.ExitOutcome{
			Stdout:   err.Error(),
			ExitCode: 1,
		}
	}

	existingObj, _ := objRes.Unwrap()

	_, _ = fmt.Fprint(os.Stdout, "enter secret data: ")
	secretRes := store.ReadSecretFromStdin()
	_, _ = fmt.Fprintln(os.Stdout)
	if secretRes.IsErr() {
		_, err := secretRes.Unwrap()
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: %v\n", err),
			ExitCode: 1,
		}
	}
	secretPlaintext, _ := secretRes.Unwrap()

	var encryptedRes fp.Result[types.Ciphertext] = fp.FlatMap(
		store.ReadPassphraseFile(inv.Env.KeyFile),
		func(pass types.Passphrase) fp.Result[types.Ciphertext] {
			return crypto.EncryptAES256CBC(secretPlaintext, pass)
		},
	)
	if encryptedRes.IsErr() {
		_, err := encryptedRes.Unwrap()
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: %v\n", err),
			ExitCode: 1,
		}
	}
	encryptedVal, _ := encryptedRes.Unwrap()

	var updatedObj map[string]any = store.PutSecret(existingObj, secretName, encryptedVal)
	var serializedRes fp.Result[[]byte] = store.SerializeJSONObject(updatedObj)
	if serializedRes.IsErr() {
		_, err := serializedRes.Unwrap()
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: %v\n", err),
			ExitCode: 1,
		}
	}
	serializedBytes, _ := serializedRes.Unwrap()

	if err := store.WriteAtomicTemp(string(secretsFile), serializedBytes); err != nil {
		return types.ExitOutcome{
			Stderr:   fmt.Sprintf("ERROR: %v\n", err),
			ExitCode: 1,
		}
	}

	return types.ExitOutcome{
		Stdout:   fmt.Sprintf("added '%s' to %s\n", secretName, secretsFile),
		ExitCode: 0,
	}
}
