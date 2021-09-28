// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"context"
	"github.com/xfali/neve-auth/user"
)

type Authenticator interface {
	AuthenticateToken(ctx context.Context, token string) (*user.UserInfo, error)
}
