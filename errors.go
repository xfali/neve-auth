// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package neveauth

import (
	"fmt"
)

var (
	malformedTokenError         = NewErr("malformed token")
	decodeTokenError            = NewErr("decoding token failed: %v")
	issuerVerifyError           = NewErr("issuer verify failed")
	tokenVerifyError            = NewErr("oidc: token verify failed: %v")
	parseClaimsError            = NewErr("oidc: parse claims failed: %v")
	indentingIdTokenClaimsError = NewErr("error indenting ID token claims: %v")
)

type Error struct {
	err string
}

func NewErr(s string) *Error {
	return &Error{
		err: s,
	}
}

func (e *Error) Error() string {
	return e.err
}

func (e *Error) V(o ...interface{}) *Error {
	return &Error{
		err: fmt.Sprintf(e.err, o...),
	}
}
