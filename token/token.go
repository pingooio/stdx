package token

import (
	"crypto/sha256"
	"errors"
	"strings"

	"github.com/pingooio/stdx/base32"
	"github.com/pingooio/stdx/crypto"
	"github.com/pingooio/stdx/guid"
)

const (
	SecretSize = crypto.KeySize256
	HashSize   = crypto.KeySize256
)

var (
	ErrTokenIsNotValid = errors.New("token is not valid")
	ErrDataIsTooLong   = errors.New("data is too long")
)

type Token struct {
	id     guid.GUID
	secret []byte
	hash   []byte
	str    string
	prefix string
}

func New(prefix string) (token Token) {
	secret := newSecret()

	return newToken(prefix, guid.NewRandom(), secret)
}

func NewWithID(prefix string, id guid.GUID) (token Token) {
	secret := newSecret()

	return newToken(prefix, id, secret)
}

// func NewWithSecret(secret []byte) (token Token, err error) {
// 	return new("", secret)
// }

// func NewWithPrefix(prefix string) (token Token, err error) {
// 	secret, err := newSecret()
// 	if err != nil {
// 		return
// 	}
// 	return newToken(prefix, secret)
// }

func newSecret() (secret []byte) {
	secret = crypto.RandBytes(SecretSize)
	return
}

func newToken(prefix string, id guid.GUID, secret []byte) (token Token) {
	idBytes, _ := id.MarshalBinary()

	hash := generateHash(idBytes, secret)

	data := append(idBytes, secret...)
	str := base32.EncodeToString(data)
	str = prefix + str

	token = Token{
		id,
		secret,
		hash,
		str,
		prefix,
	}
	return
}

func (token *Token) String() string {
	return token.str
}

func (token *Token) ID() guid.GUID {
	return token.id
}

func (token *Token) Secret() []byte {
	return token.secret
}

func (token *Token) Hash() []byte {
	return token.hash
}

func Parse(prefix, input string) (token Token, err error) {
	var tokenBytes []byte

	token.str = input

	if prefix != "" {
		if !strings.HasPrefix(input, prefix) {
			err = ErrTokenIsNotValid
			return
		}
		input = strings.TrimPrefix(input, prefix)
		token.prefix = prefix
	}

	tokenBytes, err = base32.DecodeString(input)
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}

	if len(tokenBytes) != guid.Size+SecretSize {
		err = ErrTokenIsNotValid
		return
	}

	tokenIDBytes := tokenBytes[:guid.Size]
	token.secret = tokenBytes[guid.Size:]

	token.id, err = guid.FromBytes(tokenIDBytes)
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}

	token.hash = generateHash(tokenIDBytes, token.secret)

	return
}

// FromIdAndHash creates a new token from an ID and a Hash
// it means that the token needs to be refreshed with `Refresh` before being able to use it
// as we din't have the secret, and thus cannot convert it to a valid string
func FromIdAndHash(prefix string, id guid.GUID, hash []byte) (token Token, err error) {
	if len(hash) != HashSize {
		err = ErrTokenIsNotValid
		return
	}

	token = Token{
		id:     id,
		secret: nil,
		hash:   hash,
		str:    "",
		prefix: prefix,
	}

	return
}

func (token *Token) Verify(hash []byte) (err error) {
	// in case we need to update hash size later
	// if len(hash) == OldHashSize {
	// token.hash = crypto.DeriveKeyFromKey(secret, idBytes, OldHashSize)
	// ..
	// }

	if !crypto.ConstantTimeCompare(hash, token.hash) {
		err = ErrTokenIsNotValid
	}
	return
}

func (token *Token) Refresh() {
	idBytes, _ := token.id.MarshalBinary()
	token.secret = newSecret()

	token.hash = generateHash(idBytes, token.secret)

	data := append(idBytes, token.secret...)
	str := base32.EncodeToString(data)
	token.str = token.prefix + str

	return
}

func generateHash(tokenID, secret []byte) (hash []byte) {
	hasher := sha256.New()
	hasher.Write(tokenID)
	hasher.Write(secret)
	hash = hasher.Sum(nil)
	return
}
