// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package config

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/xfali/fig"
	"github.com/xfali/neve-auth/auth"
	"github.com/xfali/neve-auth/authorizer"
	"github.com/xfali/neve-auth/filter"
	oidc2 "github.com/xfali/neve-auth/oidc"
	"github.com/xfali/neve-core/bean"
)

type FilterConfig struct{}

func NewFilterConfig() *FilterConfig {
	return &FilterConfig{}
}

func (c *FilterConfig) Init(ctx *oidc2.OidcContext, enforcer *casbin.Enforcer, conf fig.Properties, container bean.Container) error {
	verifier := ctx.GetVerifier(context.TODO())
	h := filter.NewDefaultHandler(conf, auth.NewTokenReader(), oidc2.NewAuthenticator(verifier))
	tf := filter.NewTokenFilter(h, h)
	err := container.Register(tf)
	if err != nil {
		return err
	}

	af := filter.NewAuthorizeFilter(h, authorizer.NewAuthorizer(enforcer), filter.NewUnauthHandler())
	err = container.Register(af)
	if err != nil {
		return err
	}

	return nil
}
