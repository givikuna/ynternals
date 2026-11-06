package crypto

import (
	"bytes"
	"errors"

	"github.com/givikuna/ynternals/internal/fp"
)

var ErrInvalidPadding error = errors.New("invalid PKCS#7 padding")

// PKCS7Pad pads data to a multiple of blockSize according to RFC 5652.
func PKCS7Pad(data []byte, blockSize int) []byte {
	var paddingLength int = blockSize - (len(data) % blockSize)
	var padBytes []byte = bytes.Repeat([]byte{byte(paddingLength)}, paddingLength)

	return append(append([]byte(nil), data...), padBytes...)
}

func PKCS7Unpad(data []byte, blockSize int) fp.Result[[]byte] {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return fp.Err[[]byte](ErrInvalidPadding)
	}

	var paddingLength int = int(data[len(data)-1])

	if paddingLength == 0 || paddingLength > blockSize || paddingLength > len(data) {
		return fp.Err[[]byte](ErrInvalidPadding)
	}

	var padSlice []byte = data[len(data)-paddingLength:]
	var valid bool = fp.Every(padSlice, func(b byte) bool {
		return int(b) == paddingLength
	})

	return fp.Ternary(
		valid,
		fp.Ok(append([]byte(nil), data[:len(data)-paddingLength]...)),
		fp.Err[[]byte](ErrInvalidPadding),
	)
}
