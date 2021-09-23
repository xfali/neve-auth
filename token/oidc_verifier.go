// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package token

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
)

type OidcVerifier struct {
	v *oidc.IDTokenVerifier
}

func NewOidcVerifier(v *oidc.IDTokenVerifier) *OidcVerifier {
	return &OidcVerifier{
		v: v,
	}
}

func (v *OidcVerifier) Verify(ctx context.Context, token string) (interface{}, error) {
	t, err := v.v.Verify(ctx, token)
	return t, err
}
