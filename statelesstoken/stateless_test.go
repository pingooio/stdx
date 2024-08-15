package statelesstoken_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/pingooio/stdx/crypto"
	"github.com/pingooio/stdx/statelesstoken"
	"github.com/pingooio/stdx/uuid"
)

func TestParseStateless(t *testing.T) {
	key := crypto.NewAEADKey()

	for i := 0; i < 2000; i += 1 {
		data := strconv.Itoa(i)
		newToken, err := statelesstoken.New(key, uuid.NewV4(), time.Now().Add(24*time.Hour), data)
		if err != nil {
			t.Errorf("Generating stateless token: %v", err)
		}
		parsedToken, err := statelesstoken.ParseStateless(newToken.String())
		if err != nil {
			t.Errorf("parsing  statelesstoken: %v", err)
		}

		if parsedToken.Version() != newToken.Version() {
			t.Errorf("token.Version (%v) != parsedToken.Version (%v)", newToken.Version(), parsedToken.Version())
		}

		if !parsedToken.ID().Equal(newToken.ID()) {
			t.Errorf("token.ID (%v) != parsedToken.ID (%v)", newToken.ID().String(), parsedToken.ID().String())
		}

		if parsedToken.String() != newToken.String() {
			t.Errorf("token.String() (%v) != parsedToken.String() (%v)", newToken.String(), parsedToken.String())
		}

		if parsedToken.Data() != newToken.Data() {
			t.Errorf("token.Data() (%v) != parsedToken.Data() (%v)", newToken.Data(), parsedToken.Data())
		}

		if parsedToken.Data() != data {
			t.Errorf("token.Data() (%v) != data (%v)", newToken.Data(), data)
		}
	}
}

func TestVerifyStateless(t *testing.T) {
	wrongKey := crypto.NewAEADKey()

	for i := 0; i < 2000; i += 1 {
		key := crypto.NewAEADKey()

		data := strconv.Itoa(i)
		newToken, err := statelesstoken.New(key, uuid.NewV4(), time.Now().Add(24*time.Hour), data)
		if err != nil {
			t.Errorf("Generating stateless token: %v", err)
		}

		err = newToken.Verify(key)
		if err != nil {
			t.Errorf("Verifying stateless token: %v", err)
		}

		err = newToken.Verify(wrongKey)
		if err == nil {
			t.Errorf("Accepting wrong key: %v", err)
		}
	}
}
