package validation

import (
	"encoding/json"
	"errors"

	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/types"
)

var (
	ErrInvalidJSON     = errors.New("not valid JSON")
	ErrNotJSONObject   = errors.New("must be a JSON object")
	ErrValuesNotString = errors.New("all values must be strings")
)

type pair struct {
	k string
	v any
}

func ParseRawJSON(raw []byte) fp.Result[any] {
	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return fp.Err[any](ErrInvalidJSON)
	}
	return fp.Ok(parsed)
}

func RequireJSONObject(parsed any) fp.Result[map[string]any] {
	obj, ok := parsed.(map[string]any)
	return fp.Ternary(
		ok && obj != nil,
		fp.Ok(obj),
		fp.Err[map[string]any](ErrNotJSONObject),
	)
}

func RequireStringValues(obj map[string]any) fp.Result[types.SecretMap] {
	pairs := extractPairs(obj)

	allStrings := fp.Every(pairs, func(p pair) bool {
		_, ok := p.v.(string)
		return ok
	})

	if !allStrings {
		return fp.Err[types.SecretMap](ErrValuesNotString)
	}

	secretMap := fp.Fold(pairs, make(types.SecretMap, len(pairs)), func(acc types.SecretMap, p pair) types.SecretMap {
		acc[types.SecretName(p.k)] = types.Ciphertext(p.v.(string))
		return acc
	})

	return fp.Ok(secretMap)
}

func extractPairs(obj map[string]any) []pair {
	out := make([]pair, 0, len(obj))
	for k, v := range obj {
		out = append(out, pair{k: k, v: v})
	}
	return out
}

func ValidateObjectOnly(raw []byte) fp.Result[map[string]any] {
	return fp.FlatMap(ParseRawJSON(raw), RequireJSONObject)
}

func ValidateStrictSecretMap(raw []byte) fp.Result[types.SecretMap] {
	return fp.FlatMap(ValidateObjectOnly(raw), RequireStringValues)
}
