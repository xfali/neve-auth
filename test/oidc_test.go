// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/neve-auth"
	"github.com/xfali/neve-core"
	"github.com/xfali/neve-core/processor"
	"github.com/xfali/neve-utils/neverror"
	"github.com/xfali/neve-web/gineve"
	"testing"
)

func TestOIDC(t *testing.T) {
	//m := auth.NewOidcLoginMgr()
	//m.Refresh(nil)
}

func TestRouter(t *testing.T) {
	app := neve.NewFileConfigApplication("assets/config-test.yaml")
	neverror.PanicError(app.RegisterBean(processor.NewValueProcessor()))
	neverror.PanicError(app.RegisterBean(gineve.NewProcessor()))
	neverror.PanicError(app.RegisterBean(neveauth.NewDexProcessor()))
	neverror.PanicError(app.Run())
}
