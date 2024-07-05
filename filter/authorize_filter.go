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
	"github.com/xfali/neve-auth/authorizer"
	"github.com/xfali/neve-auth/errcode"
	"github.com/xfali/xlog"
)

type GinAuthorizeFilter struct {
	logger     xlog.Logger
	attrParser AttributeParser
	authorizer authorizer.Authorizer
	fail       FailHandler
}

func NewGinAuthorizeFilter(attrParser AttributeParser, authorizer authorizer.Authorizer, fail FailHandler) *GinAuthorizeFilter {
	ret := &GinAuthorizeFilter{
		logger:     xlog.GetLogger(),
		attrParser: attrParser,
		authorizer: authorizer,
		fail:       fail,
	}
	return ret
}

func (f *GinAuthorizeFilter) Authorize(ctx *gin.Context) {
	if isIgnoreAuth(ctx) {
		ctx.Next()
		return
	}
	attr, auth, err := f.attrParser.ParseAttribute(ctx.Request)
	if err == nil {
		if !auth {
			ctx.Next()
		} else {
			if attr != nil {
				ret := f.authorizer.Authorize(ctx.Request.Context(), attr)
				if ret.Decision().IsDeny() {
					f.logger.Warnf("authorize: deny, err : %v\n", ret.Error())
					f.fail.OnFailed(errcode.DenyError, ctx.Writer, ctx.Request)
					ctx.Abort()
				} else {
					ctx.Next()
				}
			} else {
				f.fail.OnFailed(errcode.AttributeParserError, ctx.Writer, ctx.Request)
				ctx.Abort()
			}
		}
	} else {
		f.fail.OnFailed(err, ctx.Writer, ctx.Request)
		ctx.Abort()
	}
}
