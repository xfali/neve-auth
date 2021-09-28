// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/xlog"
	"net/http"
)

type TokenParser interface {
	ParseToken(ctx *gin.Context) (bool, error)
}

type TokenFilter struct {
	logger xlog.Logger
	parse  TokenParser
}

func NewTokenFilter(parse TokenParser) *TokenFilter {
	ret := &TokenFilter{
		logger: xlog.GetLogger(),
		parse:  parse,
	}
	return ret
}

func (f *TokenFilter) FilterHandler(ctx *gin.Context) {
	auth, err := f.parse.ParseToken(ctx)
	if err == nil {
		if !auth {
			ctx.Next()
		} else {
		}
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}
