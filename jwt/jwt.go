package jwt

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Algorithm string
type Type string

const (
	AlgorithmHS256 Algorithm = "HS256"
	AlgorithmHS512 Algorithm = "HAS512"
)

const (
	TypeJWT Type = "JWT"
)

var (
	ErrTokenIsNotValid     = errors.New("The token is not valid")
	ErrSignatureIsNotValid = errors.New("Signature is not valid")
	ErrTokenHasExpired     = errors.New("The token has expired")
	ErrAlgorithmIsNotValid = fmt.Errorf("Algorithm is not valid. Valid algorithms values are: [%s, %s]", AlgorithmHS256, AlgorithmHS512)
)

type Provider struct {
	signingSecretKey []byte
	algorithm        Algorithm
	verifyingKeys    [][]byte
}

type header struct {
	Algorithm Algorithm `json:"alg"`
	Type      Type      `json:"typ"`
}

// registered claim names from https://www.rfc-editor.org/rfc/rfc7519#section-4.1
type reservedClaims struct {
	ExpirationTime int64 `json:"exp,omitempty"`
	NotBefore      int64 `json:"nbf,omitempty"`
}

type NewProviderOptions struct {
	VerifyingKeys [][]byte
}

func NewProvider(signingSecretKey []byte, algorithm Algorithm, options *NewProviderOptions) (provider *Provider, err error) {
	if len(signingSecretKey) < 32 {
		err = errors.New("jwt: secretKey is too short. Min length: 32 bytes")
		return
	}

	if algorithm != AlgorithmHS256 && algorithm != AlgorithmHS512 {
		err = ErrAlgorithmIsNotValid
		return
	}

	defaultOptions := defaultNewProviderOptions()
	if options == nil {
		options = defaultOptions
	} else {
		if options.VerifyingKeys == nil {
			options.VerifyingKeys = defaultOptions.VerifyingKeys
		}
	}

	provider = &Provider{
		signingSecretKey: signingSecretKey,
		algorithm:        algorithm,
		verifyingKeys:    options.VerifyingKeys,
	}
	return
}

func defaultNewProviderOptions() *NewProviderOptions {
	return &NewProviderOptions{
		VerifyingKeys: [][]byte{},
	}
}

type TokenOptions struct {
	ExpirationTime *time.Time
	NotBefore      *time.Time
}

func (provider *Provider) IssueToken(data any, options *TokenOptions) (token string, err error) {
	tokenBuffer := bytes.NewBuffer(make([]byte, 0, 100))

	header := header{Algorithm: provider.algorithm, Type: TypeJWT}
	headerJson, err := json.Marshal(header)
	if err != nil {
		err = fmt.Errorf("jwt: encoding the header to JSON: %w", err)
		return
	}
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJson)
	tokenBuffer.WriteString(encodedHeader)
	tokenBuffer.WriteString(".")

	var claimsJson []byte
	if options != nil && (options.ExpirationTime != nil || options.NotBefore != nil) {
		var dataJson []byte
		var reservedClaims = reservedClaims{}

		if options.ExpirationTime != nil {
			reservedClaims.ExpirationTime = options.ExpirationTime.Unix()
			if reservedClaims.ExpirationTime < 1 {
				err = fmt.Errorf("jwt: ExpirationTime should not be < 1")
				return
			}
		}
		if options.NotBefore != nil {
			reservedClaims.NotBefore = options.NotBefore.Unix()
			if reservedClaims.NotBefore < 1 {
				err = fmt.Errorf("jwt: NotBefore should not be < 1")
				return
			}
		}

		claimsJson, err = json.Marshal(reservedClaims)
		if err != nil {
			err = fmt.Errorf("jwt: encoding claims to JSON: %w", err)
			return
		}
		dataJson, err = json.Marshal(data)
		if err != nil {
			err = fmt.Errorf("jwt: encoding claims to JSON: %w", err)
			return
		}
		if string(dataJson) != "{}" {
			dataJson[0] = ','
			claimsJson = append(claimsJson[:len(claimsJson)-1], dataJson...)
		}
	} else {
		claimsJson, err = json.Marshal(data)
		if err != nil {
			err = fmt.Errorf("jwt: encoding claims to JSON: %w", err)
			return
		}
		if err != nil {
			err = fmt.Errorf("jwt: encoding claims to JSON: %w", err)
			return
		}
	}

	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJson)
	tokenBuffer.WriteString(encodedClaims)

	var rawSignature []byte
	switch provider.algorithm {
	case AlgorithmHS256:
		rawSignature = signTokenHMAC(sha256.New, provider.signingSecretKey, tokenBuffer.Bytes())
	case AlgorithmHS512:
		rawSignature = signTokenHMAC(sha512.New, provider.signingSecretKey, tokenBuffer.Bytes())
	default:
		err = ErrAlgorithmIsNotValid
		return
	}
	encodedSignature := base64.RawURLEncoding.EncodeToString(rawSignature)
	tokenBuffer.WriteString(".")
	tokenBuffer.WriteString(encodedSignature)

	token = tokenBuffer.String()

	return
}

func (provider *Provider) VerifyToken(token string, data any) (err error) {
	if strings.Count(token, ".") != 2 {
		err = ErrTokenIsNotValid
		return
	}

	// Signature
	signatureStart := strings.LastIndexByte(token, '.')
	encodedSignature := token[signatureStart+1:]
	signature, err := base64.RawURLEncoding.DecodeString(encodedSignature)
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}

	encodedHeaderAndClaims := token[:signatureStart]

	switch provider.algorithm {
	case AlgorithmHS256:
		err = verifyTokenHMAC(sha256.New, provider.signingSecretKey, signature, []byte(encodedHeaderAndClaims))
	case AlgorithmHS512:
		err = verifyTokenHMAC(sha512.New, provider.signingSecretKey, signature, []byte(encodedHeaderAndClaims))
	default:
		err = ErrTokenIsNotValid
	}
	if err != nil {
		return
	}

	// Header
	var header header
	headerEnd := strings.IndexByte(token, '.')
	encodedHeader := token[:headerEnd]
	headerJson, err := base64.RawURLEncoding.DecodeString(encodedHeader)
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}
	err = json.Unmarshal(headerJson, &header)
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}

	if header.Algorithm != provider.algorithm || header.Type != TypeJWT {
		err = ErrTokenIsNotValid
		return
	}

	// Reserved Claims
	encodedClaims := token[headerEnd+1 : signatureStart]
	claimsJson, err := base64.RawURLEncoding.DecodeString(encodedClaims)
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}

	var reservedClaims reservedClaims
	err = json.Unmarshal(claimsJson, &reservedClaims)
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}

	now := time.Now().Unix()
	if reservedClaims.ExpirationTime != 0 {
		if now > reservedClaims.ExpirationTime {
			err = ErrTokenHasExpired
			return
		}
	}
	if reservedClaims.NotBefore != 0 {
		if now < reservedClaims.NotBefore {
			err = ErrTokenIsNotValid
			return
		}
	}

	err = json.Unmarshal(claimsJson, data)
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}

	return
}
