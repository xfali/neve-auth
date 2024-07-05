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

package test

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth"
	"github.com/xfali/neve-auth/filter"
	"github.com/xfali/neve-core"
	"github.com/xfali/neve-core/processor"
	"github.com/xfali/neve-utils/neverror"
	"github.com/xfali/neve-web/gineve"
	"github.com/xfali/neve-web/gineve/midware/loghttp"
	"github.com/xfali/neve-web/result"
	"net/http"
	"testing"
)

type webBean struct {
	V          string             //`fig:"Log.Level"`
	HttpLogger loghttp.HttpLogger `inject:""`
	filter.PermissionHandlerHolder
}

func (b *webBean) HttpRoutes(engine gin.IRouter) {
	engine.GET("test", b.HttpLogger.LogHttp(), b.RequirePermission("userinfo", "read"), func(context *gin.Context) {
		context.JSON(http.StatusOK, result.Ok(b.V))
	})
}

func TestOIDC(t *testing.T) {
	//m := auth.NewOidcLoginMgr()
	//m.Refresh(nil)
}

func TestRouter(t *testing.T) {
	app := neve.NewFileConfigApplication("assets/config-test.yaml")
	neverror.PanicError(app.RegisterBean(processor.NewValueProcessor()))
	neverror.PanicError(app.RegisterBean(gineve.NewProcessor()))
	neverror.PanicError(app.RegisterBean(neveauth.NewDexProcessor()))
	neverror.PanicError(app.RegisterBean(&webBean{}))
	neverror.PanicError(app.Run())
}
