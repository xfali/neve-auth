// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/auth"
	"github.com/xfali/xlog"
)

type AuthRouter struct {
	logger xlog.Logger

	RedirectUrl string `fig:"neve.auth.web.router.redirect"`
	CallbackUrl string `fig:"neve.auth.web.router.callback"`

	redirectHandler auth.RedirectHandler
	callbackHandler auth.CallbackHandler
	refreshHandler  auth.RefreshHandler
}

func (r *AuthRouter) BeanAfterSet() error {
	if r.RedirectUrl == "" {
		err := errors.New("redirect url is empty")
		r.logger.Errorln(err)
		return err
	}

	if r.CallbackUrl == "" {
		err := errors.New("callback url is empty")
		r.logger.Errorln(err)
		return err
	}
	return nil
}

func (r *AuthRouter) HttpRoutes(engine gin.IRouter) {
	engine.POST(r.RedirectUrl, r.redirect)
	engine.GET(r.CallbackUrl, r.callback)
	engine.POST(r.CallbackUrl, r.refresh)
}

func (r *AuthRouter) redirect(ctx *gin.Context) {
	r.redirectHandler.Redirect(ctx)
}

func (r *AuthRouter) callback(ctx *gin.Context) {
	r.callbackHandler.Callback(ctx)
}

func (r *AuthRouter) refresh(ctx *gin.Context) {
	r.refreshHandler.Refresh(ctx)
}
