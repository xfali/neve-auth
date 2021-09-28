// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package errcode

import (
	"fmt"
)

var (
	MalformedTokenError         = NewErr("malformed token")
	DecodeTokenError            = NewErr("decoding token failed: %v")
	IssuerVerifyError           = NewErr("issuer verify failed")
	TokenVerifyError            = NewErr("oidc: token verify failed: %v")
	ParseClaimsError            = NewErr("oidc: parse claims failed: %v")
	IndentingIdTokenClaimsError = NewErr("error indenting ID token claims: %v")
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
