package store

import (
	"encoding/json"
	"slices"

	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/types"
)

func PutSecret(existing map[string]any, name types.SecretName, value types.Ciphertext) map[string]any {
	next := cloneAnyMap(existing)
	next[string(name)] = string(value)
	return next
}

func DeleteSecret(existing map[string]any, name types.SecretName) map[string]any {
	next := cloneAnyMap(existing)
	delete(next, string(name))
	return next
}

func SortedKeys(obj map[string]any) []string {
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func SortedEntries(secrets types.SecretMap) []types.SecretEntry {
	names := make([]string, 0, len(secrets))
	for k := range secrets {
		names = append(names, string(k))
	}
	slices.Sort(names)

	return fp.Map(names, func(rawName string) types.SecretEntry {
		name := types.SecretName(rawName)
		return types.SecretEntry{
			Name:  name,
			Value: secrets[name],
		}
	})
}

func SerializeJSONObject(obj map[string]any) fp.Result[[]byte] {
	encoded, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fp.Err[[]byte](err)
	}
	return fp.Ok(append(encoded, '\n'))
}

func cloneAnyMap(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
