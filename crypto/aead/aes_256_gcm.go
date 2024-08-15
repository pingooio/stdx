package aead

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/pingooio/stdx/crypto"
)

const (
	// Aes256GcmKeySize is the size of the key used by the AES-256-GCM AEAD, in bytes.
	Aes256GcmKeySize = crypto.Size256

	// Aes256GcmNonceSize is the size of the nonce used with the AES-256-GCM
	// variant of this AEAD, in bytes.
	Aes256GcmNonceSize = crypto.Size96
)

// NewAes256GcmKey generates a new random secret key.
func NewAes256GcmKey() []byte {
	return crypto.RandBytes(Aes256GcmKeySize)
}

// NewAes256GcmKey generates a new random nonce.
func NewAes256GcmNonce() []byte {
	return crypto.RandBytes(Aes256GcmNonceSize)
}

// NewAEAD returns a AES-256-GCM AEAD that uses the given 256-bit key.
func NewAes256Gcm(key []byte) (aeadCipher cipher.AEAD, err error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	aeadCipher, err = cipher.NewGCM(blockCipher)
	if err != nil {
		return
	}
	return
}
