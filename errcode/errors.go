/*
 * Copyright (C) 2019-2024, Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

	AttributeParserError = NewErr("parse attribute failed")

	DenyError = NewErr("deny")
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
