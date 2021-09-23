// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package neveauth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/coreos/go-oidc/v3/oidc"
	oidc2 "github.com/xfali/neve-auth/oidc"
	"github.com/xfali/neve-auth/token"
	"strings"
)

type oidcAuthenticator struct {
	oidcCtx  *oidc2.OidcContext
	verifier token.Verifier
	scopes   []string
}

func CreateAuthenticator(oidcCtx *oidc2.OidcContext) (*oidcAuthenticator, error) {
	ret := &oidcAuthenticator{
		oidcCtx: oidcCtx,
	}

	ret.verifier = oidcCtx.GetVerifier(context.Background())
	return ret, nil
}

func (a *oidcAuthenticator) AuthenticateToken(ctx context.Context, token string) (*UserInfo, error) {
	//err := a.VerifyIssuer(token)
	//if err != nil {
	//	return err
	//}
	t, err := a.verifier.Verify(ctx, token)
	if err != nil {
		return nil, tokenVerifyError.V(err)
	}
	idToken := t.(*oidc.IDToken)

	var claims json.RawMessage
	if err := idToken.Claims(&claims); err != nil {
		return nil, parseClaimsError.V(err)
	}

	buff := new(bytes.Buffer)
	if err := json.Indent(buff, []byte(claims), "", "  "); err != nil {
		return nil, indentingIdTokenClaimsError.V(err)
	}

	return &UserInfo{}, nil
}

func (a *oidcAuthenticator) VerifyIssuer(token string) error {
	issuer, err := parseIssuer(token)
	if err != nil {
		return err
	}
	if issuer != a.oidcCtx.IssuerURL {
		return issuerVerifyError
	}
	return nil
}

func parseIssuer(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", malformedTokenError
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", decodeTokenError.V(err)
	}
	return string(payload), nil
}
