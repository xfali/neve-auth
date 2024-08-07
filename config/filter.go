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

package config

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/xfali/fig"
	"github.com/xfali/neve-auth/auth"
	"github.com/xfali/neve-auth/authorizer"
	"github.com/xfali/neve-auth/constants"
	"github.com/xfali/neve-auth/filter"
	oidc2 "github.com/xfali/neve-auth/oidc"
	"github.com/xfali/neve-core/bean"
)

type FilterConfig struct{}

func NewFilterConfig() *FilterConfig {
	return &FilterConfig{}
}

func (c *FilterConfig) Init(ctx *oidc2.OidcContext, redirectHandler auth.RedirectHandler, enforcer *casbin.Enforcer, conf fig.Properties, container bean.Container) error {
	verifier := ctx.GetVerifier(context.TODO())
	matcher := filter.NewPathMatcher(conf)
	ignoreDexPaths(conf, matcher)
	h := filter.NewDefaultHandler(matcher, redirectHandler, auth.NewTokenReader())
	tf := filter.NewTokenFilter(h, oidc2.NewAuthenticator(verifier), h)
	err := container.Register(tf)
	if err != nil {
		return err
	}

	af := filter.NewGinAuthorizeFilter(h, authorizer.NewAuthorizer(enforcer), filter.NewUnauthHandler())
	ph := filter.NewPermissionHandler(af)

	err = container.RegisterByName(filter.PermissionHandlerName, ph)
	if err != nil {
		return err
	}
	return nil
}

func ignoreDexPaths(conf fig.Properties, matcher filter.PathMatcher) {
	url := conf.Get(constants.RouterRedirectKey, "")
	if url != "" {
		matcher.Exclude(url)
	}

	url = conf.Get(constants.RouterUserInfoKey, "")
	if url != "" {
		matcher.Exclude(url)
	}

	url = conf.Get(constants.RouterCallbackKey, "")
	if url != "" {
		matcher.Exclude(url)
	}
}
