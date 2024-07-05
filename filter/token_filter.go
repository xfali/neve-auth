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

package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/attribute"
	"github.com/xfali/xlog"
	"net/http"
)

type TokenParser interface {
	ParseToken(req *http.Request) (string, bool, error)
}

type TokenFilter struct {
	logger        xlog.Logger
	parse         TokenParser
	fail          FailHandler
	authenticator Authenticator
}

func NewTokenFilter(parse TokenParser, authenticator Authenticator, fail FailHandler) *TokenFilter {
	ret := &TokenFilter{
		logger:        xlog.GetLogger(),
		parse:         parse,
		fail:          fail,
		authenticator: authenticator,
	}
	return ret
}

func (f *TokenFilter) FilterHandler(ctx *gin.Context) {
	token, auth, err := f.parse.ParseToken(ctx.Request)
	if err == nil {
		if auth {
			user, err := f.authenticator.AuthenticateToken(ctx.Request.Context(), token)
			if err != nil {
				f.fail.OnFailed(err, ctx.Writer, ctx.Request)
				ctx.Abort()
			}
			ctx.Request = attribute.WithUser(ctx.Request, user)
		} else {
			setIgnoreAuth(ctx)
		}
		ctx.Next()
	} else {
		f.fail.OnFailed(err, ctx.Writer, ctx.Request)
		ctx.Abort()
	}
}
