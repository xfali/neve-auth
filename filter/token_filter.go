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
	ParseToken(req *http.Request) (bool, error)
}

type TokenFilter struct {
	logger xlog.Logger
	parse  TokenParser
	fail   FailHandler
}

func NewTokenFilter(parse TokenParser, fail FailHandler) *TokenFilter {
	ret := &TokenFilter{
		logger: xlog.GetLogger(),
		parse:  parse,
		fail:   fail,
	}
	return ret
}

func (f *TokenFilter) FilterHandler(ctx *gin.Context) {
	auth, err := f.parse.ParseToken(ctx.Request)
	if err == nil {
		if !auth {
			ctx.Next()
		}
	} else {
		f.fail.OnFailed(err, ctx.Writer, ctx.Request)
		ctx.Abort()
	}
}
