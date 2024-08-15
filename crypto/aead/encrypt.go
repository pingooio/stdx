package aead

import (
	"crypto/cipher"
	"errors"

	"github.com/pingooio/stdx/crypto"
)

// Encrypt is an helper function to symetrically encrypt a piece of data using AES-256-GCM
// returning the nonce separatly
func EncryptWithNonce(key, plaintext, additionalData []byte) (ciphertext, nonce []byte, err error) {
	nonce = NewAes256GcmNonce()

	cipher, err := NewAes256Gcm(key)
	if err != nil {
		return
	}

	ciphertext = cipher.Seal(nil, nonce, plaintext, additionalData)
	return
}

// DecryptWithNonce is an helper function to symetrically  decrypt a piece of data using AES-256-GCM
// taking the nonce as a separate piece of input
func DecryptWithNonce(key, nonce, ciphertext, additionalData []byte) (plaintext []byte, err error) {
	cipher, err := NewAes256Gcm(key)
	if err != nil {
		return
	}

	plaintext, err = cipher.Open(nil, nonce, ciphertext, additionalData)
	return
}

// Encrypt is an helper function to symetrically encrypt a piece of data using the given cipher
// the nonce is prepended to the ciphertext in the returned buffer
func Encrypt(cipher cipher.AEAD, key, plaintext, additionalData []byte) (ciphertext []byte) {
	nonce := crypto.RandBytes(uint64(cipher.NonceSize()))

	ciphertext = cipher.Seal(nil, nonce, plaintext, additionalData)
	ciphertext = append(nonce, ciphertext...)
	return
}

// DecryptWithNonce is an helper function to symetrically decrypt a piece of data using the provided cipher
// The nonce should be at the begining of the ciphertext
func Decrypt(cipher cipher.AEAD, key, ciphertext, additionalData []byte) (plaintext []byte, err error) {
	nonceSize := cipher.NonceSize()

	if len(ciphertext) < nonceSize {
		err = errors.New("crypto.Decrypt: len(ciphertext) < NonceSize")
		return
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	plaintext, err = cipher.Open(nil, nonce, ciphertext, additionalData)
	return
}
