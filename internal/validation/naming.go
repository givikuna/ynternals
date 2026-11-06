package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/types"
)

var validNamePattern *regexp.Regexp = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func IsValidSecretName(raw string) bool {
	return validNamePattern.MatchString(raw) && !strings.Contains(raw, "..")
}

func ParseSecretName(raw string) fp.Result[types.SecretName] {
	return fp.Ternary(
		IsValidSecretName(raw),
		fp.Ok(types.SecretName(raw)),
		fp.Err[types.SecretName](fmt.Errorf("invalid secret name: %s", raw)),
	)
}
