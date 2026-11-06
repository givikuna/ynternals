package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/givikuna/ynternals/internal/fp"
	"github.com/givikuna/ynternals/internal/types"
)

var (
	opensslMagic         []byte = []byte("Salted__")
	ErrInvalidCiphertext error  = errors.New("invalid OpenSSL ciphertext format")
)

func EncryptAES256CBC(plaintext types.Plaintext, pass types.Passphrase) fp.Result[types.Ciphertext] {
	var salt types.Salt

	if _, err := rand.Read(salt[:]); err != nil {
		return fp.Err[types.Ciphertext](err)
	}

	return EncryptWithSalt(plaintext, pass, salt)
}

func EncryptWithSalt(plaintext types.Plaintext, pass types.Passphrase, salt types.Salt) fp.Result[types.Ciphertext] {
	var material types.DerivedKeyMaterial = DeriveOpenSSLKeyAndIV(pass, salt)

	block, err := aes.NewCipher(material.Key[:])
	if err != nil {
		return fp.Err[types.Ciphertext](err)
	}

	var padded []byte = PKCS7Pad(plaintext, aes.BlockSize)
	var encrypted []byte = make([]byte, len(padded))

	cipher.NewCBCEncrypter(block, material.IV[:]).CryptBlocks(encrypted, padded)

	var payload []byte = append(append(append([]byte(nil), opensslMagic...), salt[:]...), encrypted...)
	var encoded string = base64.StdEncoding.EncodeToString(payload)

	return fp.Ok(types.Ciphertext(encoded))
}

func DecryptAES256CBC(ciphertext types.Ciphertext, pass types.Passphrase) fp.Result[types.Plaintext] {
	var cleanedStr string = strings.Map(func(r rune) rune {
		return fp.Ternary(r == '\n' || r == '\r' || r == ' ' || r == '\t', -1, r)
	}, string(ciphertext))

	raw, err := base64.StdEncoding.DecodeString(cleanedStr)
	if err != nil {
		return fp.Err[types.Plaintext](err)
	}

	var headerSize int = len(opensslMagic) + len(types.Salt{})
	if len(raw) < headerSize+aes.BlockSize || !bytes.Equal(raw[:len(opensslMagic)], opensslMagic) {
		return fp.Err[types.Plaintext](ErrInvalidCiphertext)
	}

	var salt types.Salt
	copy(salt[:], raw[len(opensslMagic):headerSize])
	var encBytes []byte = raw[headerSize:]

	if len(encBytes)%aes.BlockSize != 0 {
		return fp.Err[types.Plaintext](ErrInvalidCiphertext)
	}

	var material types.DerivedKeyMaterial = DeriveOpenSSLKeyAndIV(pass, salt)
	block, err := aes.NewCipher(material.Key[:])
	if err != nil {
		return fp.Err[types.Plaintext](err)
	}

	var decrypted []byte = make([]byte, len(encBytes))
	cipher.NewCBCDecrypter(block, material.IV[:]).CryptBlocks(decrypted, encBytes)

	return fp.MapResult(PKCS7Unpad(decrypted, aes.BlockSize), func(b []byte) types.Plaintext {
		return types.Plaintext(b)
	})
}
