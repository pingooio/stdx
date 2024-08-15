package jwt

import (
	"crypto/hmac"
	"crypto/subtle"
	"hash"
)

func signTokenHMAC(hashFunction func() hash.Hash, secret, encodedHeaderAndClaims []byte) (signatureRaw []byte) {
	hmac := hmac.New(hashFunction, secret)
	hmac.Write(encodedHeaderAndClaims)
	signatureRaw = hmac.Sum(nil)

	return
}

func verifyTokenHMAC(hashFunction func() hash.Hash, secret, signature, encodedHeaderAndClaims []byte) (err error) {
	hmac := hmac.New(hashFunction, secret)
	hmac.Write(encodedHeaderAndClaims)
	hmacHash := hmac.Sum(nil)

	if subtle.ConstantTimeCompare(hmacHash, signature) != 1 {
		err = ErrSignatureIsNotValid
		return
	}
	return
}
