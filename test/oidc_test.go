// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth"
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
	V          string //`fig:"Log.Level"`
	HttpLogger loghttp.HttpLogger `inject:""`
}

func (b *webBean) HttpRoutes(engine gin.IRouter) {
	engine.GET("test", b.HttpLogger.LogHttp(), neveauth.Gin.RequirePermission("user", "read"), func(context *gin.Context) {
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
