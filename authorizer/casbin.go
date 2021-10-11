// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package authorizer

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/xfali/neve-auth/attribute"
)

type defaultAuthorizer struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizer(enforcer *casbin.Enforcer) *defaultAuthorizer {
	ret := &defaultAuthorizer{
		enforcer: enforcer,
	}
	return ret
}

func (a *defaultAuthorizer) Authorize(ctx context.Context, attr attribute.Attribute) Result {
	subject := attr.GetUserInfo().Username
	project := attr.GetUserInfo().ProjectID
	resource := attr.GetResource()
	action := attr.GetAction()

	ok, err := a.enforcer.Enforce(subject, project, resource, action)
	if err != nil {
		return Deny(err)
	}
	if !ok {
		return Deny(nil)
	}
	return Allow()
}

type MergedAuthorizer struct {
	authorizers []Authorizer
}

func MergeAuthorizer(authorizers ...Authorizer) Authorizer {
	return &MergedAuthorizer{
		authorizers: authorizers,
	}
}

func (a *MergedAuthorizer) Authorize(ctx context.Context, attr attribute.Attribute) (ret Result) {
	for _, auth := range a.authorizers {
		ret = auth.Authorize(ctx, attr)
		switch ret.Decision() {
		case DecisionAllow, DecisionDeny:
			return ret
		}
	}
	return
}
