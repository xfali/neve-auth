/*
 * Copyright (C) 2019-2024, Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
