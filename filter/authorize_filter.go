// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/authorizer"
	"github.com/xfali/neve-auth/errcode"
	"github.com/xfali/xlog"
)

type AuthorizeFilter struct {
	logger     xlog.Logger
	attrParser AttributeParser
	authorizer authorizer.Authorizer
	fail       FailHandler
}

func NewAuthorizeFilter(attrParser AttributeParser, authorizer authorizer.Authorizer, fail FailHandler) *AuthorizeFilter {
	ret := &AuthorizeFilter{
		logger:     xlog.GetLogger(),
		attrParser: attrParser,
		authorizer: authorizer,
		fail:       fail,
	}
	return ret
}

func (f *AuthorizeFilter) FilterHandler(ctx *gin.Context) {
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
