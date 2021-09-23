// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package authorizer

import (
	"context"
	"github.com/casbin/casbin/v2"
)

type defaultAuthorizer struct {
	enforcer *casbin.SyncedEnforcer
}

func NewAuthorizer(enforcer *casbin.SyncedEnforcer) *defaultAuthorizer {
	ret := &defaultAuthorizer{
		enforcer: enforcer,
	}
	return ret
}

func (a *defaultAuthorizer) Authorize(ctx context.Context, attr Attribute) Result {
	subject := attr.UserInfo().Username

}
