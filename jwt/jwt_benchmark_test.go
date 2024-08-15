package jwt_test

import (
	"testing"

	"github.com/pingooio/stdx/jwt"
)

func BenchmarkJwt(b *testing.B) {
	jwtProvider, _ := jwt.NewProvider(TEST_INSECURE_SECRET, jwt.AlgorithmHS256, nil)
	type payload struct {
		Message string
	}
	claims := payload{Message: "Hello World"}
	var parsedClaims payload

	b.Run("issue", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = jwtProvider.IssueToken(claims, nil)
		}
	})

	token, _ := jwtProvider.IssueToken(claims, nil)

	b.Run("verify", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = jwtProvider.VerifyToken(token, &parsedClaims)
		}
	})
}
