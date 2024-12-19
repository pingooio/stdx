package bchacha20blake3

import (
	"crypto/cipher"
	"crypto/subtle"
	"encoding/binary"
	"errors"

	"github.com/pingooio/stdx/crypto/blake3"
	"github.com/pingooio/stdx/crypto/chacha20"
	// "golang.org/x/crypto/chacha20"
)

const (
	KeySize   = 32
	NonceSize = 32
	TagSize   = 32

	encryptionKeyContext    = "BChaCha20-BLAKE3 2023-12-31 23:59:59:999 encryption key"
	athenticationKeyContext = "BChaCha20-BLAKE3 2024-01-01 00:00:00:000 authentication key"
)

var (
	ErrOpen = errors.New("bchacha20blake3: error decrypting ciphertext")
)

type BChaCha20Blake3 struct {
	key [KeySize]byte
}

// ensure that BChaCha20Blake3 implements `cipher.AEAD` interface at build time
var _ cipher.AEAD = (*BChaCha20Blake3)(nil)

func New(key []byte) (*BChaCha20Blake3, error) {
	ret := new(BChaCha20Blake3)
	copy(ret.key[:], key)
	return ret, nil
}

func (*BChaCha20Blake3) NonceSize() int {
	return NonceSize
}

func (*BChaCha20Blake3) Overhead() int {
	return TagSize
}

func (x *BChaCha20Blake3) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	var encryptionKey [32]byte
	var authenticationKey [32]byte

	deriveKey(encryptionKey[:0], x.key[:], encryptionKeyContext, nil)
	deriveKey(authenticationKey[:0], x.key[:], athenticationKeyContext, nonce)

	ret, out := sliceForAppend(dst, len(plaintext)+TagSize)
	ciphertext, tag := out[:len(plaintext)], out[len(plaintext):]

	chacha20Cipher, _ := chacha20.New(encryptionKey[:], nonce[24:32])
	chacha20Cipher.XORKeyStream(ciphertext, plaintext)

	// _ = tag
	macHasher := blake3.New(32, authenticationKey[:])
	macHasher.Write(additionalData)
	// macHasher.Write(nonce)
	macHasher.Write(ciphertext)
	writeUint64LittleEndian(macHasher, uint64(len(additionalData)))
	// writeUint64(macHasher, uint64(len(nonce)))
	writeUint64LittleEndian(macHasher, uint64(len(ciphertext)))
	macHasher.Sum(tag[:0])

	return ret
}

func (x *BChaCha20Blake3) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	var encryptionKey [32]byte
	var authenticationKey [32]byte

	deriveKey(encryptionKey[:0], x.key[:], encryptionKeyContext, nil)

	deriveKey(authenticationKey[:0], x.key[:], athenticationKeyContext, nonce)

	tag := ciphertext[len(ciphertext)-TagSize:]
	ciphertext = ciphertext[:len(ciphertext)-TagSize]

	chacha20Cipher, _ := chacha20.New(encryptionKey[:], nonce[24:32])

	var computedTag [TagSize]byte
	macHasher := blake3.New(32, authenticationKey[:])
	macHasher.Write(additionalData)
	// macHasher.Write(nonce)
	macHasher.Write(ciphertext)
	writeUint64LittleEndian(macHasher, uint64(len(additionalData)))
	// writeUint64(macHasher, uint64(len(nonce)))
	writeUint64LittleEndian(macHasher, uint64(len(ciphertext)))
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

func deriveKey(out, parentKey []byte, context string, nonce []byte) {
	var keyMaterial [12 + KeySize]byte

	copy(keyMaterial[0:12], nonce)
	copy(keyMaterial[12:44], parentKey)
	// binary.LittleEndian.PutUint64(keyMaterial[44:52], uint64(len(nonce)))
	// binary.LittleEndian.PutUint64(keyMaterial[52:60], uint64(len(parentKey)))

	// blake3x.DeriveKey(out, context, keyMaterial[:])

	blake3.DeriveKey(out, context, keyMaterial[:])

	// hasher := blake3.NewDeriveKey(context)
	// hasher.Write(nonce)
	// hasher.Write(parentKey)
	// writeUint64(hasher, uint64(len(nonce)))
	// writeUint64(hasher, uint64(len(parentKey)))
	// hasher.Sum(out)
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

func writeUint64LittleEndian(p *blake3.Hasher, n uint64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], n)
	p.Write(buf[:])
}
