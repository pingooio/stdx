package chacha20blake3

import (
	"crypto/cipher"
	"errors"

	"github.com/pingooio/stdx/crypto/chacha"
)

const (
	NonceSizeX = 24
)

var (
	ErrBadNonceXLength = errors.New("chacha20blake3: bad nonce length for XChaCha20Blake3. 24 bytes required")
)

type XChaCha20Blake3 struct {
	key [KeySize]byte
}

// ensure that XChaCha20Blake3 implements `cipher.AEAD` interface at build time
var _ cipher.AEAD = (*XChaCha20Blake3)(nil)

func NewX(key []byte) (*XChaCha20Blake3, error) {
	if len(key) != KeySize {
		return nil, ErrBadKeyLength
	}

	ret := new(XChaCha20Blake3)
	copy(ret.key[:], key)
	return ret, nil
}

func (*XChaCha20Blake3) NonceSize() int {
	return NonceSizeX
}

func (*XChaCha20Blake3) Overhead() int {
	return TagSize
}

func (x *XChaCha20Blake3) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	if len(nonce) != NonceSizeX {
		panic(ErrBadNonceXLength)
	}

	chaChaKey, _ := chacha.HChaCha20(x.key[:], nonce[0:16])
	chacha20Cipher, _ := New(chaChaKey)

	// chaChaNonce := make([]byte, NonceSize)
	// copy(chaChaNonce[4:12], nonce[16:24])

	return chacha20Cipher.Seal(dst, nonce[16:24], plaintext, additionalData)
}

func (x *XChaCha20Blake3) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	if len(nonce) != NonceSizeX {
		panic(ErrBadNonceXLength)
	}

	chaChaKey, _ := chacha.HChaCha20(x.key[:], nonce[0:16])
	chacha20Cipher, _ := New(chaChaKey)

	return chacha20Cipher.Open(dst, nonce[16:24], ciphertext, additionalData)
}
