package token_test

import (
	"strings"
	"testing"

	"github.com/pingooio/stdx/crypto"
	"github.com/pingooio/stdx/token"
	"github.com/pingooio/stdx/uuid"
)

var emptyUUID = make([]byte, uuid.Size)
var emptySecret = make([]byte, token.SecretSize)
var emptyHash = make([]byte, token.HashSize)

func TestNew(t *testing.T) {
	newToken := token.New("")

	tokenIDByte, _ := newToken.ID().MarshalBinary()
	if crypto.ConstantTimeCompare(tokenIDByte, emptyUUID) {
		t.Error("token.ID is empty")
	}

	if len(newToken.Secret()) != token.SecretSize {
		t.Errorf("Bad secret size. Expected: %d | Got: %d", token.SecretSize, len(newToken.Secret()))
	}

	if crypto.ConstantTimeCompare(newToken.Secret(), emptySecret) {
		t.Error("token.Secret is empty")
	}

	if len(newToken.Hash()) != token.HashSize {
		t.Errorf("Bad hash size. Expected: %d | Got: %d", token.HashSize, len(newToken.Hash()))
	}

	if crypto.ConstantTimeCompare(newToken.Hash(), emptyHash) {
		t.Error("token.Hash is empty")
	}
}

func TestNewWithPrefix(t *testing.T) {
	prefix := "test_"
	newToken := token.New(prefix)

	if !strings.HasPrefix(newToken.String(), prefix) {
		t.Errorf("Token doesn't have prefix. Expected: %s | Got: %s", prefix, newToken.String())
	}
}

func TestParse(t *testing.T) {
	var err error
	for i := 0; i < 1000; i += 1 {
		newToken := token.New("")
		_, err = token.Parse("", newToken.String())
		if err != nil {
			t.Errorf("parsing token: %v", err)
		}
	}
}

func TestVerify(t *testing.T) {
	var err error

	for i := 0; i < 1000; i += 1 {
		newToken := token.New("")

		if err = newToken.Verify(newToken.Hash()); err != nil {
			t.Errorf("verifying otken. expected: nil | got: %v", err)
		}
	}

	newToken := token.New("")
	if err = newToken.Verify(nil); err == nil {
		t.Errorf("verifying token against null.  expected: %v | got: %v", token.ErrTokenIsNotValid, err)
	}
	if err = newToken.Verify([]byte{}); err == nil {
		t.Errorf("verifying token against empty slice.  expected: %v | got: %v", token.ErrTokenIsNotValid, err)
	}
}
