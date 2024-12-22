package bchacha20blake3

import (
	"bytes"
	"crypto/rand"
	"errors"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	var key [KeySize]byte
	var nonce [NonceSize]byte

	originalPlaintext := []byte("Hello World")
	additionalData := []byte("!")

	rand.Read(key[:])
	rand.Read(nonce[:])

	cipher, _ := New(key[:])
	ciphertext := cipher.Seal(nil, nonce[:], originalPlaintext, additionalData)

	decryptedPlaintext, err := cipher.Open(nil, nonce[:], ciphertext, additionalData)
	if err != nil {
		t.Errorf("decrypting message: %s", err)
		return
	}

	if !bytes.Equal(decryptedPlaintext, originalPlaintext) {
		t.Errorf("original message (%s) != decrypted message (%s)", string(originalPlaintext), string(decryptedPlaintext))
		return
	}
}

func TestAdditionalData(t *testing.T) {
	var key [KeySize]byte
	var nonce [NonceSize]byte

	originalPlaintext := []byte("Hello World")
	additionalData := []byte("!")

	rand.Read(key[:])
	rand.Read(nonce[:])

	cipher, _ := New(key[:])
	ciphertext := cipher.Seal(nil, nonce[:], originalPlaintext, additionalData)

	_, err := cipher.Open(nil, nonce[:], ciphertext, []byte{})
	if !errors.Is(err, ErrOpen) {
		t.Errorf("expected error (%s) | got (%s)", ErrOpen, err)
		return
	}
}

func BenchmarkZeroize(b *testing.B) {
	var key [32]byte

	b.Run("zeroizeKey", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(32)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			zeroizeKey(key)
		}
	})

	b.Run("zeroize", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(32)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			zeroize(key[:])
		}
	})
}

// Deprecated: Use zeroize instead.
func zeroizeKey(input [32]byte) {
	var zeros [32]byte
	copy(input[:], zeros[:])
}
