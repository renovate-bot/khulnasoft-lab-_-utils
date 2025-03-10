package oidctest_test

import (
	"context"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-jose/go-jose/v3/jwt"

	"github.com/khulnasoft-lab/utils/oidctest"
)

func TestNewIssuer(t *testing.T) {
	ctx := context.Background()

	signer, iss := oidctest.NewIssuer(t)

	token, err := jwt.Signed(signer).Claims(jwt.Claims{
		Issuer:   iss,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Expiry:   jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		Subject:  "test-subject",
		Audience: jwt.Audience{"test-audience"},
	}).CompactSerialize()
	if err != nil {
		t.Fatalf("CompactSerialize() = %v", err)
	}

	// Verify the token is valid.
	provider, err := oidc.NewProvider(ctx, iss)
	if err != nil {
		t.Errorf("constructing %q provider: %v", iss, err)
	}

	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
	if _, err := verifier.Verify(ctx, token); err != nil {
		t.Errorf("verifying token: %v", err)
	}
}
