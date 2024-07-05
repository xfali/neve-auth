/*
 * Copyright (C) 2024, Xiongfa Li.
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

package digest

import (
	"context"
	auth "github.com/abbot/go-http-auth"
	"github.com/xfali/neve-auth/user"
	"net/http"
	"time"
)

type digestAuthenticator struct {
	auth     *auth.DigestAuth
	Verifier user.UserPasswordVerifier
}

func NewDigestAuthenticator(realm string) *digestAuthenticator {
	ret := &digestAuthenticator{}
	ret.auth = auth.NewDigestAuthenticator(realm, ret.secret)
	ret.auth.PlainTextSecrets = true
	return ret
}

func (h *digestAuthenticator) secret(user, realm string) string {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	pw, err := h.Verifier.QueryPassword(ctx, user)
	if err == nil {
		return pw
	}
	return ""
}

func (h *digestAuthenticator) Wrap(f http.HandlerFunc) http.HandlerFunc {
	return h.auth.Wrap(func(writer http.ResponseWriter, request *auth.AuthenticatedRequest) {
		f(writer, &request.Request)
	})
}
