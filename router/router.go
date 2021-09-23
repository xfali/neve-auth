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

	RedirectUrl string `fig:"neve.auth.router.redirect"`
	CallbackUrl string `fig:"neve.auth.router.callback"`
	UserInfoUrl string `fig:"neve.auth.router.userinfo"`

	redirectHandler auth.RedirectHandler
	callbackHandler auth.CallbackHandler
	refreshHandler  auth.RefreshHandler
	userInfoHandler auth.UserInfoHandler
}

func NewAuthRouter(
	redirectHandler auth.RedirectHandler,
	callbackHandler auth.CallbackHandler,
	refreshHandler auth.RefreshHandler,
	userInfoHandler auth.UserInfoHandler) *AuthRouter {
	return &AuthRouter{
		logger:          xlog.GetLogger(),
		redirectHandler: redirectHandler,
		callbackHandler: callbackHandler,
		refreshHandler:  refreshHandler,
		userInfoHandler: userInfoHandler,
	}
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
	engine.GET(r.UserInfoUrl, r.info)
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

func (r *AuthRouter) info(ctx *gin.Context) {
	r.userInfoHandler.GetUserInfo(ctx)
}
