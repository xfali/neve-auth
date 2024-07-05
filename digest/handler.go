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
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/user"
)

type AuthenticateHandler interface {
	DigestAuthenticate() gin.HandlerFunc
}

type authenticateHandler struct {
	auth     *digestAuthenticator
	Verifier user.UserPasswordVerifier `inject:""`

	Realm string `fig:"neve.auth.digest.realm"`
}

func NewAuthenticateHandler() *authenticateHandler {
	return &authenticateHandler{}
}

func (o *authenticateHandler) BeanAfterSet() error {
	if o.Realm == "" {
		o.Realm = "neve-auth-digest-realm"
	}
	o.auth = NewDigestAuthenticator(o.Realm)
	if o.Verifier == nil {
		panic("UserPasswordVerifier is nil ")
	}
	o.auth.Verifier = o.Verifier
	return nil
}

func (o *authenticateHandler) DigestAuthenticate() gin.HandlerFunc {
	return func(context *gin.Context) {
		if username, authinfo := o.auth.auth.CheckAuth(context.Request); username == "" {
			o.auth.auth.RequireAuth(context.Writer, context.Request)
			context.Abort()
		} else {
			//ar := &auth.AuthenticatedRequest{Request: *context.Request, Username: username}
			if authinfo != nil {
				context.Writer.Header().Set(o.auth.auth.Headers.V().AuthInfo, *authinfo)
			}
			context.Next()
		}
	}
}
