package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"hash"

	"github.com/givikuna/ynternals/internal/types"
)

const (
	opensslPBKDF2Iterations = 10000
	aesKeySize              = 32
	aesIVSize               = 16
	derivedTotalSize        = aesKeySize + aesIVSize
)

func DeriveOpenSSLKeyAndIV(pass types.Passphrase, salt types.Salt) types.DerivedKeyMaterial {
	var dk []byte = pbkdf2SHA256(pass, salt[:], opensslPBKDF2Iterations, derivedTotalSize)
	var material types.DerivedKeyMaterial

	copy(material.Key[:], dk[:aesKeySize])
	copy(material.IV[:], dk[aesKeySize:derivedTotalSize])

	return material
}

func pbkdf2SHA256(password, salt []byte, iterations, keyLen int) []byte {
	var numBlocks int = (keyLen + sha256.Size - 1) / sha256.Size
	var full []byte = deriveBlocks(password, salt, iterations, 1, numBlocks)

	return full[:keyLen]
}

func deriveBlocks(password, salt []byte, iterations, currentBlock, totalBlocks int) []byte {
	if currentBlock > totalBlocks {
		return nil
	}

	var block []byte = deriveSingleBlock(password, salt, iterations, currentBlock)

	return append(block, deriveBlocks(password, salt, iterations, currentBlock+1, totalBlocks)...)
}

func deriveSingleBlock(password, salt []byte, iterations, blockIndex int) []byte {
	var intBuf [4]byte
	binary.BigEndian.PutUint32(intBuf[:], uint32(blockIndex))

	var u1 []byte = prfSHA256(password, append(append([]byte(nil), salt...), intBuf[:]...))

	return iteratePRF(password, u1, u1, iterations-1)
}

func iteratePRF(password, prevU, acc []byte, remaining int) []byte {
	if remaining <= 0 {
		return acc
	}

	var nextU []byte = prfSHA256(password, prevU)

	return iteratePRF(password, nextU, xorSlices(acc, nextU), remaining-1)
}

func prfSHA256(key, data []byte) []byte {
	var mac hash.Hash = hmac.New(sha256.New, key)
	mac.Write(data)

	return mac.Sum(nil)
}

func xorSlices(a, b []byte) []byte {
	if len(a) == 0 {
		return nil
	}

	return append([]byte{a[0] ^ b[0]}, xorSlices(a[1:], b[1:])...)
}
