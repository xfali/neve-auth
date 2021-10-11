// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/gin-gonic/gin"
)

const (
	ctxIngoreKey = "_neve_auth_ignore_"
)

func setIgnoreAuth(ctx *gin.Context) {
	ctx.Set(ctxIngoreKey, true)
}

func isIgnoreAuth(ctx *gin.Context) bool {
	v, ok := ctx.Get(ctxIngoreKey)
	if ok {
		return v.(bool)
	}
	return false
}
