// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package oidc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/xfali/neve-auth/errcode"
	"github.com/xfali/neve-auth/token"
	"github.com/xfali/neve-auth/user"
	"strings"
)

type oidcAuthenticator struct {
	verifier token.Verifier
	scopes   []string
}

func NewAuthenticator(verifier token.Verifier) *oidcAuthenticator {
	ret := &oidcAuthenticator{
		verifier: verifier,
	}
	return ret
}

func (a *oidcAuthenticator) AuthenticateToken(ctx context.Context, token string) (*user.UserInfo, error) {
	//err := a.VerifyIssuer(token)
	//if err != nil {
	//	return err
	//}
	t, err := a.verifier.Verify(ctx, token)
	if err != nil {
		return nil, errcode.TokenVerifyError.V(err)
	}
	idToken := t.(*oidc.IDToken)

	var claims json.RawMessage
	if err := idToken.Claims(&claims); err != nil {
		return nil, errcode.ParseClaimsError.V(err)
	}

	buff := new(bytes.Buffer)
	if err := json.Indent(buff, []byte(claims), "", "  "); err != nil {
		return nil, errcode.IndentingIdTokenClaimsError.V(err)
	}

	return &user.UserInfo{}, nil
}

func (a *oidcAuthenticator) VerifyIssuer(token string) error {
	//issuer, err := parseIssuer(token)
	//if err != nil {
	//	return err
	//}
	//if issuer != a.oidcCtx.IssuerURL {
	//	return errcode.IssuerVerifyError
	//}
	return nil
}

func parseIssuer(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errcode.MalformedTokenError
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errcode.DecodeTokenError.V(err)
	}
	return string(payload), nil
}
