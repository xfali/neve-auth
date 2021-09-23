// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import "github.com/gin-gonic/gin"

type PathFilter interface {
	ShouldFilter(path string) bool
}

type OidcFilter struct {
	pf PathFilter
}

func NewOidcFilter() *OidcFilter {
	ret := &OidcFilter{}
	return ret
}

func (f *OidcFilter) FilterHandler(ctx *gin.Context) {
	if f.pf.ShouldFilter(ctx.Request.RequestURI) {

	} else {
		ctx.Next()
	}
}
