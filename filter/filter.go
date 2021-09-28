// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/authorizer"
	"github.com/xfali/xlog"
	"net/http"
)

type AuthFilter struct {
	logger     xlog.Logger
	attrParser AttributeParser
	authorizer authorizer.Authorizer
}

func NewAuthFilter(attrParser AttributeParser, authorizer authorizer.Authorizer) *AuthFilter {
	ret := &AuthFilter{
		logger:     xlog.GetLogger(),
		attrParser: attrParser,
		authorizer: authorizer,
	}
	return ret
}

func (f *AuthFilter) FilterHandler(ctx *gin.Context) {
	attr, err := f.attrParser.ParseAttribute(ctx)
	if err == nil {
		if attr != nil {
			ret := f.authorizer.Authorize(ctx.Request.Context(), attr)
			if ret.Decision().IsDeny() {
				ctx.AbortWithStatus(http.StatusUnauthorized)
			} else {
				ctx.Next()
			}
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	} else {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
}
