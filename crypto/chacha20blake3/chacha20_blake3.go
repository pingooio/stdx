package chacha20blake3

import (
	"crypto/cipher"
	"crypto/subtle"
	"encoding/binary"
	"errors"

	"github.com/pingooio/stdx/crypto/blake3"
	// "golang.org/x/crypto/chacha20"
	"github.com/pingooio/stdx/crypto/chacha20"
)

const (
	KeySize   = 32
	NonceSize = 8
	TagSize   = 32
)

var (
	ErrOpen           = errors.New("chacha20blake3: error decrypting ciphertext")
	ErrBadKeyLength   = errors.New("chacha20blake3: bad key length for ChaCha20Blake3. 32 bytes required")
	ErrBadNonceLength = errors.New("chacha20blake3: bad nonce length for ChaCha20Blake3. 8 bytes required")
)

type ChaCha20Blake3 struct {
	key [KeySize]byte
}

// ensure that ChaCha20Blake3 implements `cipher.AEAD` interface at build time
var _ cipher.AEAD = (*ChaCha20Blake3)(nil)

func New(key []byte) (*ChaCha20Blake3, error) {
	if len(key) != KeySize {
		return nil, ErrBadKeyLength
	}

	var ret ChaCha20Blake3
	copy(ret.key[:], key)

	return &ret, nil
}

func (*ChaCha20Blake3) NonceSize() int {
	return NonceSize
}

func (*ChaCha20Blake3) Overhead() int {
	return TagSize
}

func (cipher *ChaCha20Blake3) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	var authenticationKey [32]byte

	if len(nonce) != NonceSize {
		panic(ErrBadNonceLength)
	}

	ret, out := sliceForAppend(dst, len(plaintext)+TagSize)
	ciphertext, tag := out[:len(plaintext)], out[len(plaintext):]

	// first get a new ChaCha20 cipher instance
	chacha20Cipher, _ := chacha20.New(cipher.key[:], nonce)
	// var ietfNonce [12]byte
	// copy(ietfNonce[:], nonce)
	// chacha20Cipher, _ := chacha20.NewUnauthenticatedCipher(cipher.key[:], ietfNonce[:])

	// then perform the KDF step to get the authentication key and increase the ChaCha20 counter
	chacha20Cipher.XORKeyStream(authenticationKey[:], authenticationKey[:])
	chacha20Cipher.SetCounter(1)

	// then encrypt the plaintext
	chacha20Cipher.XORKeyStream(ciphertext, plaintext)

	// and finally MAC the AAD + ciphertext with the authentication key
	macHasher := blake3.New(32, authenticationKey[:])
	macHasher.Write(additionalData)
	macHasher.Write(ciphertext)
	writeUint64(macHasher, uint64(len(additionalData)))
	writeUint64(macHasher, uint64(len(ciphertext)))
	macHasher.Sum(tag[:0])

	return ret
}

func (cipher *ChaCha20Blake3) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	var authenticationKey [32]byte

	if len(nonce) != NonceSize {
		panic(ErrBadNonceLength)
	}

	tag := ciphertext[len(ciphertext)-TagSize:]
	ciphertext = ciphertext[:len(ciphertext)-TagSize]

	chacha20Cipher, _ := chacha20.New(cipher.key[:], nonce)
	// var ietfNonce [12]byte
	// copy(ietfNonce[:], nonce)
	// chacha20Cipher, _ := chacha20.NewUnauthenticatedCipher(cipher.key[:], ietfNonce[:])

	chacha20Cipher.XORKeyStream(authenticationKey[:], authenticationKey[:])
	chacha20Cipher.SetCounter(1)

	var computedTag [TagSize]byte
	macHasher := blake3.New(32, authenticationKey[:])
	macHasher.Write(additionalData)
	macHasher.Write(ciphertext)
	writeUint64(macHasher, uint64(len(additionalData)))
	writeUint64(macHasher, uint64(len(ciphertext)))
	macHasher.Sum(computedTag[:0])

	ret, plaintext := sliceForAppend(dst, len(ciphertext))

	if subtle.ConstantTimeCompare(computedTag[:], tag) != 1 {
		// for i := range plaintext {
		// 	plaintext[i] = 0
		// }
		return nil, ErrOpen
	}

	chacha20Cipher.XORKeyStream(plaintext, ciphertext)

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
