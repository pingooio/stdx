package schacha20blake3

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"

	"github.com/pingooio/stdx/crypto/blake3"
	"github.com/pingooio/stdx/crypto/chacha20"
	"github.com/pingooio/stdx/crypto/chacha20blake3"
	// "golang.org/x/crypto/chacha20"
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
	// var subKey [32]byte
	// copy(subKey[:], nonce[0:16])

	// chacha20Kdf1, _ := chacha20.New(x.key[:], nonce[16:24])
	// // chacha20Kdf1, _ := chacha20.NewUnauthenticatedCipher(x.key[:], nonce[16:28])
	// chacha20Kdf1.XORKeyStream(subKey[:], subKey[:])

	// chacha20Blake3Cipher, _ := chacha20blake3.New(subKey[:])
	// return chacha20Blake3Cipher.Seal(dst, nonce[24:32], plaintext, additionalData)

	// var authenticationKey [32]byte
	var encryptionKey [32]byte
	copy(encryptionKey[:], nonce[0:8])

	var authenticationKey [32]byte
	copy(authenticationKey[:], nonce[8:16])

	chacha20Kdf, _ := chacha20.New(x.key[:], nonce[16:24])
	// chacha20Kdf, _ := chacha20.NewUnauthenticatedCipher(x.key[:], nonce[0:12])
	chacha20Kdf.XORKeyStream(authenticationKey[:], authenticationKey[:])
	// we encrypted / hashed 32 bytes. Now we skip the 32 bytes remaining from the 64 bytes block
	chacha20Kdf.SetCounter(1)
	chacha20Kdf.XORKeyStream(encryptionKey[:], encryptionKey[:])
	// chacha20Kdf.SetCounter(2)

	ret, out := sliceForAppend(dst, len(plaintext)+TagSize)
	ciphertext, tag := out[:len(plaintext)], out[len(plaintext):]

	chacha20Cipher, _ := chacha20.New(encryptionKey[:], nonce[24:32])
	// chacha20Cipher, _ := chacha20.NewUnauthenticatedCipher(x.key[:], nonce[0:12])
	chacha20Cipher.XORKeyStream(ciphertext, plaintext)

	macHasher := blake3.New(32, authenticationKey[:])
	macHasher.Write(additionalData)
	macHasher.Write(ciphertext)
	writeUint64(macHasher, uint64(len(additionalData)))
	writeUint64(macHasher, uint64(len(ciphertext)))
	macHasher.Sum(tag[:0])

	zeroize(encryptionKey[:])
	zeroize(authenticationKey[:])

	return ret
}

func (x *SChaCha20Blake3) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	var subKey [32]byte
	copy(subKey[0:16], nonce[8:24])

	chacha20Kdf, _ := chacha20.New(x.key[:], nonce[0:8])
	// chacha20Kdf, _ := chacha20.NewUnauthenticatedCipher(x.key[:], nonce[0:12])
	chacha20Kdf.XORKeyStream(subKey[:], subKey[:])

	chacha20Blake3Cipher, _ := chacha20blake3.New(subKey[:])
	ret, err := chacha20Blake3Cipher.Open(dst, nonce[24:32], ciphertext, additionalData)
	if err != nil {
		return nil, ErrOpen
	}
	return ret, nil
}

// sliceForAppend takes a slice and a requested number of bytes. It returns a
// slice with the contents of the given slice followed by that many bytes and a
// second slice that aliases into it and contains only the extra bytes. If the
// original slice has sufficient capacity then no allocation is performed.
func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

func writeUint64(p *blake3.Hasher, n uint64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], n)
	p.Write(buf[:])
}

func zeroize(input []byte) {
	for i := range input {
		input[i] = 0
	}
}
