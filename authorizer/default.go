// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package authorizer

import "github.com/xfali/neve-auth/user"

type defaultAttribute struct {
	user *user.UserInfo
}

func (a *defaultAttribute) UserInfo() *user.UserInfo {
	return a.user
}
