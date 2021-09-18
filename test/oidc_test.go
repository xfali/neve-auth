// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/neve-auth/auth"
	"testing"
)

func TestOIDC(t *testing.T) {
	m := auth.NewOidcLoginMgr()
	m.Refresh(nil)
}
