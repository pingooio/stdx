package jwt_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/pingooio/stdx/jwt"
)

var (
	TEST_INSECURE_SECRET = []byte("9MfG7mC+dqxXgXZGdiFsxyET1mQRBekiw7K+bIvZEGc=")
)

func TestVerifyInvalidTokens(t *testing.T) {
	jwtProvider, _ := jwt.NewProvider(TEST_INSECURE_SECRET, jwt.AlgorithmHS256, nil)
	invalidTokens := []string{
		".",
		"..",
		"...",
		"....",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..",
		".eyJtZXNzYWdlIjoiOSJ9.",
		"..pRVbZ4OC40z7qZsRj0cPpLodPmMAF7-skUqztxeK9iQ=",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZXNzYWdlIjoiOSJ9.",
		".eyJtZXNzYWdlIjoiOSJ9.pRVbZ4OC40z7qZsRj0cPpLodPmMAF7-skUqztxeK9iQ=",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..pRVbZ4OC40z7qZsRj0cPpLodPmMAF7-skUqztxeK9iQ=",
	}

	for _, invalidToken := range invalidTokens {
		var decodedPayload struct{}
		err := jwtProvider.VerifyToken(invalidToken, &decodedPayload)
		if err == nil {
			t.Errorf("The invalid token: %s is verifed as valid", invalidToken)
		}
	}
}

func TestVerifyKnownJWTs(t *testing.T) {
	jwts := []struct {
		Token     string
		Algorithm jwt.Algorithm
		Secret    []byte
	}{
		{
			Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZXNzYWdlIjoiSGVsbG8gV29ybGQifQ.B3yZ_TwtqZGk1ejyUI6QcD-eqrPN0awNEy4L4CiwvL8",
			Algorithm: jwt.AlgorithmHS256,
			Secret:    []byte("secretsecretsecretsecretsecretsecret"),
		},
	}

	for _, knowJwt := range jwts {
		var emptyStruct struct{}
		jwtProvider, _ := jwt.NewProvider(knowJwt.Secret, jwt.AlgorithmHS256, nil)
		err := jwtProvider.VerifyToken(knowJwt.Token, &emptyStruct)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestIssueAndVerify(t *testing.T) {
	type claims struct {
		Message string `json:"message"`
	}
	jwtProvider, _ := jwt.NewProvider(TEST_INSECURE_SECRET, jwt.AlgorithmHS256, nil)

	for i := 0; i < 10; i += 1 {
		payload := claims{Message: strconv.Itoa(i)}
		token, err := jwtProvider.IssueToken(payload, nil)
		if err != nil {
			t.Error(err)
		}

		var decodedPayload claims
		err = jwtProvider.VerifyToken(token, &decodedPayload)
		if err != nil {
			t.Error(err)
		}
		if decodedPayload.Message != payload.Message {
			t.Errorf("payload.message (%s) != decodedPayload.message(%s)", payload.Message, decodedPayload.Message)
		}
	}
}

func TestIssueAndVerifyInvalid(t *testing.T) {
	type claims struct {
		Message string `json:"message"`
	}
	jwtProvider, _ := jwt.NewProvider(TEST_INSECURE_SECRET, jwt.AlgorithmHS256, nil)

	for i := 0; i < 10; i += 1 {
		payload := claims{Message: strconv.Itoa(i)}
		token, err := jwtProvider.IssueToken(payload, nil)
		if err != nil {
			t.Error(err)
		}

		var decodedPayload claims

		invalidTokenWithPrefix := "x" + token
		err = jwtProvider.VerifyToken(invalidTokenWithPrefix, &decodedPayload)
		if !errors.Is(err, jwt.ErrSignatureIsNotValid) {
			t.Errorf("expected error: %s. got: %v", jwt.ErrSignatureIsNotValid, err)
		}

		invalidTokenWithInvalidPrefix := "|" + token
		err = jwtProvider.VerifyToken(invalidTokenWithInvalidPrefix, &decodedPayload)
		if !errors.Is(err, jwt.ErrSignatureIsNotValid) {
			t.Errorf("expected error: %s. got: %v", jwt.ErrSignatureIsNotValid, err)
		}

		invalidTokenWithSuffix := token + "x"
		err = jwtProvider.VerifyToken(invalidTokenWithSuffix, &decodedPayload)
		if !errors.Is(err, jwt.ErrSignatureIsNotValid) {
			t.Errorf("expected error: %s. got: %v", jwt.ErrSignatureIsNotValid, err)
		}

		invalidTokenWithInvalidSuffix := token + "|"
		err = jwtProvider.VerifyToken(invalidTokenWithInvalidSuffix, &decodedPayload)
		if !errors.Is(err, jwt.ErrTokenIsNotValid) {
			t.Errorf("expected error: %s. got: %v", jwt.ErrTokenIsNotValid, err)
		}
	}
}

func TestIssueAndVerifyExpireIntTheFuture(t *testing.T) {
	type claims struct {
		Message string `json:"message"`
	}
	jwtProvider, _ := jwt.NewProvider(TEST_INSECURE_SECRET, jwt.AlgorithmHS256, nil)

	for i := 0; i < 10; i += 1 {
		payload := claims{Message: strconv.Itoa(i)}
		expiresAt := time.Now().Add(1 * time.Hour)
		token, err := jwtProvider.IssueToken(payload, &jwt.TokenOptions{ExpirationTime: &expiresAt})
		if err != nil {
			t.Error(err)
		}

		var decodedPayload claims
		err = jwtProvider.VerifyToken(token, &decodedPayload)
		if err != nil {
			t.Error(err)
		}
		if decodedPayload.Message != payload.Message {
			t.Errorf("pyaload.message (%s) != decodedPayload.message(%s)", payload.Message, decodedPayload.Message)
		}
	}
}

func TestIssueAndVerifyExpired(t *testing.T) {
	type emptyClaims struct {
	}
	type claims struct {
		Message string `json:"message"`
	}
	jwtProvider, _ := jwt.NewProvider(TEST_INSECURE_SECRET, jwt.AlgorithmHS256, nil)

	for i := 0; i < 10; i += 1 {
		expiresAt := time.Now().Add(-1 * time.Hour)
		token, err := jwtProvider.IssueToken(emptyClaims{}, &jwt.TokenOptions{ExpirationTime: &expiresAt})
		if err != nil {
			t.Error(err)
		}

		var decodedPayload emptyClaims
		err = jwtProvider.VerifyToken(token, &decodedPayload)
		if !errors.Is(err, jwt.ErrTokenHasExpired) {
			t.Errorf("expected error: %s. got: %v", jwt.ErrTokenHasExpired, err)
		}
	}

	for i := 0; i < 10; i += 1 {
		payload := claims{Message: strconv.Itoa(i)}
		expiresAt := time.Now().Add(-1 * time.Hour)
		token, err := jwtProvider.IssueToken(payload, &jwt.TokenOptions{ExpirationTime: &expiresAt})
		if err != nil {
			t.Error(err)
		}

		var decodedPayload claims
		err = jwtProvider.VerifyToken(token, &decodedPayload)
		if !errors.Is(err, jwt.ErrTokenHasExpired) {
			t.Errorf("expected error: %s. got: %v", jwt.ErrTokenHasExpired, err)
		}
	}
}
