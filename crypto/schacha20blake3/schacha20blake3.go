package schacha20blake3

import (
	"crypto/cipher"
	"errors"

	"github.com/pingooio/stdx/crypto/chacha20"
	"github.com/pingooio/stdx/crypto/chacha20blake3"
	// "golang.org/x/crypto/chacha20"
	// "github.com/pingooio/stdx/crypto/chacha20"
)

const (
	KeySize   = 32
	NonceSize = 32
	TagSize   = 32
)

var (
	ErrOpen = errors.New("xchacha20blake3: error decrypting ciphertext")
)

type SChaCha20Blake3 struct {
	key [KeySize]byte
}

// ensure that SChaCha20Blake3 implements `cipher.AEAD` interface at build time
var _ cipher.AEAD = (*SChaCha20Blake3)(nil)

func New(key []byte) (*SChaCha20Blake3, error) {
	ret := new(SChaCha20Blake3)
	copy(ret.key[:], key)
	return ret, nil
}

func (*SChaCha20Blake3) NonceSize() int {
	return NonceSize
}

func (*SChaCha20Blake3) Overhead() int {
	return TagSize
}

func (x *SChaCha20Blake3) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	// ret, out := sliceForAppend(dst, len(plaintext)+TagSize)
	// ciphertext, tag := out[:len(plaintext)], out[len(plaintext):]

	// var authenticationKey [32]byte
	var subKey [32]byte
	copy(subKey[0:16], nonce[8:24])

	chacha20Kdf, _ := chacha20.New(x.key[:], nonce[0:8])
	chacha20Kdf.XORKeyStream(subKey[:], subKey[:])

	chacha20Blake3Cipher, _ := chacha20blake3.New(subKey[:])
	return chacha20Blake3Cipher.Seal(dst, nonce[24:32], plaintext, additionalData)

	// chacha20Cipher.(authenticationKey[:], authenticationKey[:])
	// chacha20Cipher.SetCounter(1)
	// chacha20Cipher.XORKeyStream(ciphertext, plaintext)

	// macHasher, _ := blake3.NewKeyed(authenticationKey[:])
	// macHasher.Write(additionalData)
	// macHasher.Write(ciphertext)
	// writeUint64(macHasher, uint64(len(additionalData)))
	// writeUint64(macHasher, uint64(len(ciphertext)))
	// macHasher.Sum(tag[:0])

	// return ret
}

func (x *SChaCha20Blake3) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	var subKey [32]byte
	copy(subKey[0:16], nonce[8:24])

	chacha20Kdf, _ := chacha20.New(x.key[:], nonce[0:8])
	chacha20Kdf.XORKeyStream(subKey[:], subKey[:])

	chacha20Blake3Cipher, _ := chacha20blake3.New(subKey[:])
	ret, err := chacha20Blake3Cipher.Open(dst, nonce[24:32], ciphertext, additionalData)
	if err != nil {
		return nil, ErrOpen
	}
	return ret, nil
}
