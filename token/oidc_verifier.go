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
