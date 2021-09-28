// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package authorizer

import "github.com/xfali/neve-auth/user"

type defaultAttribute struct {
	User     *user.UserInfo
	Resource string
	Action   string
	Tenant   string
}

func (a *defaultAttribute) GetUserInfo() *user.UserInfo {
	return a.User
}

func (a *defaultAttribute) GetAction() string {
	return a.Action
}

func (a *defaultAttribute) GetResource() string {
	return a.Resource
}
