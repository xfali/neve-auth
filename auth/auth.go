// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package auth

import "github.com/gin-gonic/gin"

type RedirectHandler interface {
	Redirect(ctx *gin.Context)
}

type CallbackHandler interface {
	Callback(ctx *gin.Context)
}

type RefreshHandler interface {
	Refresh(ctx *gin.Context)
}
