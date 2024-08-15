package kdf

import (
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/hkdf"
)

func HkdfSha256(secret, info, salt []byte, size int) (key []byte, err error) {
	if size < 1 {
		err = errors.New("hkdf: size is not valid")
		return
	}
	if secret == nil {
		err = errors.New("hkdf: size can't be null")
		return
	}
	if info == nil {
		err = errors.New("hkdf: info can't be null")
		return
	}

	key = make([]byte, size)
	hkdf := hkdf.New(sha256.New, secret, salt, info)
	_, err = io.ReadFull(hkdf, key)
	return
}
