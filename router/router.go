/*
 * Copyright (C) 2019-2024, Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/auth"
	"github.com/xfali/xlog"
)

type AuthRouter struct {
	logger xlog.Logger

	RedirectUrl string `fig:"neve.auth.authenticator.spec.router.redirect"`
	CallbackUrl string `fig:"neve.auth.authenticator.spec.router.callback"`
	UserInfoUrl string `fig:"neve.auth.authenticator.spec.router.userinfo"`

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
	r.redirectHandler.Redirect(ctx.Writer, ctx.Request)
}

func (r *AuthRouter) callback(ctx *gin.Context) {
	r.callbackHandler.Callback(ctx.Writer, ctx.Request)
}

func (r *AuthRouter) refresh(ctx *gin.Context) {
	r.refreshHandler.Refresh(ctx.Writer, ctx.Request)
}

func (r *AuthRouter) info(ctx *gin.Context) {
	r.userInfoHandler.GetUserInfo(ctx.Writer, ctx.Request)
}
