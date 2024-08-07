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

package neveauth

import (
	"context"
	"fmt"
	oidc2 "github.com/coreos/go-oidc/v3/oidc"
	"github.com/xfali/fig"
	"github.com/xfali/neve-auth/auth"
	"github.com/xfali/neve-auth/config"
	"github.com/xfali/neve-auth/constants"
	"github.com/xfali/neve-auth/oidc"
	"github.com/xfali/neve-auth/router"
	"github.com/xfali/neve-core/bean"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type dexOpts struct {
}

type dexProcessor struct {
	dexOpts
}

type DexOpt func(*dexOpts)

func NewDexProcessor(opts ...DexOpt) *dexProcessor {
	ret := &dexProcessor{}
	for _, opt := range opts {
		opt(&ret.dexOpts)
	}
	return ret
}

// 初始化对象处理器
func (p *dexProcessor) Init(conf fig.Properties, container bean.Container) error {
	client := getClient(conf)
	issuer := conf.Get(constants.IssuerKey, "")
	if issuer == "" {
		return fmt.Errorf("issuer value is empty, set it: %s", constants.IssuerKey)
	}

	oauthConf, err := getOAuthConfig(conf)
	if err != nil {
		return err
	}

	ctx := oidc.NewOidcContext(client, issuer, oauthConf)
	err = ctx.Init()
	if err != nil {
		return err
	}

	enforcer, err := config.NewCasbinConfig().Init(conf, container)
	if err != nil {
		return err
	}

	oidcMgr := auth.NewOidcLoginMgr(oauthConf, client)
	err = config.NewFilterConfig().Init(ctx, oidcMgr, enforcer, conf, container)
	if err != nil {
		return err
	}
	return container.Register(router.NewAuthRouter(oidcMgr, oidcMgr, oidcMgr, oidcMgr))
}

// 对象分类，判断对象是否实现某些接口，并进行相关归类。为了支持多协程处理，该方法应线程安全。
// 注意：该方法建议只做归类，具体处理使用Process，不保证Processor的实现在此方法中做了相关处理。
// 该方法在Bean Inject注入之后调用
// return: bool 是否能够处理对象， error 处理是否有错误
func (p *dexProcessor) Classify(o interface{}) (bool, error) {
	return false, nil
}

// 对已分类对象做统一处理，注意如果存在耗时操作，请使用其他协程处理。
// 该方法在Classify及BeanAfterSet后调用。
// 成功返回nil，失败返回error
func (p *dexProcessor) Process() error {
	return nil
}

func (p *dexProcessor) BeanDestroy() error {
	return nil
}

func getClient(conf fig.Properties) *http.Client {
	return http.DefaultClient
}

func getOAuthConfig(conf fig.Properties) (*oauth2.Config, error) {
	issuer := conf.Get(constants.IssuerKey, "")
	if issuer == "" {
		return nil, fmt.Errorf("issuer value is empty, set it: %s", constants.IssuerKey)
	}

	id := conf.Get(constants.ClientIdKey, "")
	if id == "" {
		return nil, fmt.Errorf("client id value is empty, set it: %s", constants.ClientIdKey)
	}

	secret := conf.Get(constants.ClientSecretKey, "")
	if secret == "" {
		return nil, fmt.Errorf("client secret value is empty, set it: %s", constants.ClientSecretKey)
	}

	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	provider, err := oidc2.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to query provider %q: %v", issuer, err)
	}

	url, err := constants.RedirectUrl(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to query redirect url %s: %v", url, err)
	}
	return &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc2.ScopeOpenID, oidc2.ScopeOfflineAccess, "profile", "email", "groups"},
		RedirectURL:  url,
	}, nil
}
