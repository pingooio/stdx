package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const (
	// AEADKeySize is the size of the key used by this AEAD, in bytes.
	AEADKeySize = KeySize256

	// AEADNonceSize is the size of the nonce used with the AES-256-GCM
	// variant of this AEAD, in bytes.
	AEADNonceSize = 12
)

// NewAEADKey generates a new random secret key.
func NewAEADKey() []byte {
	return RandBytes(AEADKeySize)
}

// NewAEADNonce generates a new random nonce.
func NewAEADNonce() []byte {
	return RandBytes(AEADNonceSize)
}

// NewAEAD returns a AES-256-GCM AEAD that uses the given 256-bit key.
func NewAEAD(key []byte) (aeadCipher cipher.AEAD, err error) {
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

// Encrypt is an helper function to symetrically encrypt a piece of data using AES-256-GCM
// returning the nonce separatly
func EncryptWithNonce(key, plaintext, additionalData []byte) (ciphertext, nonce []byte, err error) {
	nonce = NewAEADNonce()

	cipher, err := NewAEAD(key)
	if err != nil {
		return
	}
	ciphertext = cipher.Seal(nil, nonce, plaintext, additionalData)
	return
}

// DecryptWithNonce is an helper function to symetrically  decrypt a piece of data using AES-256-GCM
// taking the nonce as a separate piece of input
func DecryptWithNonce(key, nonce, ciphertext, additionalData []byte) (plaintext []byte, err error) {
	cipher, err := NewAEAD(key)
	if err != nil {
		return
	}
	plaintext, err = cipher.Open(nil, nonce, ciphertext, additionalData)
	return
}

// Encrypt is an helper function to symetrically encrypt a piece of data using AES-256-GCM
// the nonce is prepended to the ciphertext in the returned buffer
func Encrypt(key, plaintext, additionalData []byte) (ciphertext []byte, err error) {
	nonce := NewAEADNonce()

	cipher, err := NewAEAD(key)
	if err != nil {
		return
	}
	ciphertext = cipher.Seal(nil, nonce, plaintext, additionalData)
	ciphertext = append(nonce, ciphertext...)
	return
}

// DecryptWithNonce is an helper function to symetrically decrypt a piece of data using AES-256-GCM
// The nonce should be at the begining of the ciphertext
func Decrypt(key, ciphertext, additionalData []byte) (plaintext []byte, err error) {
	cipher, err := NewAEAD(key)
	if err != nil {
		return
	}

	if len(ciphertext) < AEADNonceSize {
		err = errors.New("crypto.Decrypt: len(ciphertext) < NonceSize")
		return
	}
	nonce := ciphertext[:AEADNonceSize]
	ciphertext = ciphertext[AEADNonceSize:]

	plaintext, err = cipher.Open(nil, nonce, ciphertext, additionalData)
	return
}
