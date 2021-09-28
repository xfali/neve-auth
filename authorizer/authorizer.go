// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package authorizer

import (
	"context"
	"github.com/xfali/neve-auth/user"
)

type Attribute interface {
	GetUserInfo() *user.UserInfo
	GetAction() string
	GetResource() string
}

type Authorizer interface {
	Authorize(ctx context.Context, attr Attribute) Result
}
