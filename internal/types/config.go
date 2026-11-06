package types

type SecretsFilePath string

type KeyFilePath string

type RuntimeDirPath string

const (
	DefaultSecretsFile SecretsFilePath = "secrets.json"
	DefaultKeyFile     KeyFilePath     = "/etc/nixos/symmetric.key"
	DefaultRuntimeDir  RuntimeDirPath  = "/run/ynternals"
	KeyFileEnvVar      string          = "YNTERNALS_KEYFILE"
)

type EnvConfig struct {
	EUID       int
	KeyFile    KeyFilePath
	RuntimeDir RuntimeDirPath
	Version    string
}
