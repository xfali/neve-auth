// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package attribute

import "github.com/xfali/neve-auth/user"

type DefaultAttribute struct {
	User     *user.UserInfo
	Resource string
	Action   string
	Tenant   string
}

func (a *DefaultAttribute) GetUserInfo() *user.UserInfo {
	return a.User
}

func (a *DefaultAttribute) GetAction() string {
	return a.Action
}

func (a *DefaultAttribute) GetResource() string {
	return a.Resource
}
