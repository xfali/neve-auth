// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package auth

import "github.com/gin-gonic/gin"

type TokenReader interface {
	ReadToken(ctx *gin.Context) (*Token, error)
}

type TokenWriter interface {
	WriteToken(ctx *gin.Context, token *Token) error
}

type RedirectHandler interface {
	Redirect(ctx *gin.Context)
}

type CallbackHandler interface {
	Callback(ctx *gin.Context)
}

type RefreshHandler interface {
	Refresh(ctx *gin.Context)
}

type UserInfoHandler interface {
	GetUserInfo(ctx *gin.Context)
}
