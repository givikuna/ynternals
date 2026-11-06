package types

type SecretName string

type Plaintext []byte

type Ciphertext string

type Passphrase []byte

type Salt [8]byte

type DerivedKeyMaterial struct {
	Key [32]byte
	IV  [16]byte
}

type SecretEntry struct {
	Name  SecretName
	Value Ciphertext
}

type SecretMap map[SecretName]Ciphertext
