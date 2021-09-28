// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package neveauth

import (
	"context"
	"fmt"
	oidc2 "github.com/coreos/go-oidc/v3/oidc"
	"github.com/xfali/fig"
	"github.com/xfali/neve-auth/auth"
	"github.com/xfali/neve-auth/config"
	"github.com/xfali/neve-auth/router"
	"github.com/xfali/neve-core/bean"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

const (
	issuerKey       = "neve.auth.dex.issuer"
	clientIdKey     = "neve.auth.dex.client.id"
	clientSecretKey = "neve.auth.dex.client.secret"
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
	oauthConf, err := getOAuthConfig(conf)
	if err != nil {
		return err
	}
	//err := container.Register(oidc.NewOidcContext(client, issuer, oauthConf))
	//if err != nil {
	//	return err
	//}

	err = config.NewCasbinConfig().Init(conf, container)
	if err != nil {
		return err
	}

	oidcMgr := auth.NewOidcLoginMgr(oauthConf, client)
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
	issuer := conf.Get(issuerKey, "")
	if issuer == "" {
		return nil, fmt.Errorf("issuer value is empty, set it: %s", issuerKey)
	}

	id := conf.Get(clientIdKey, "")
	if id == "" {
		return nil, fmt.Errorf("client id value is empty, set it: %s", clientIdKey)
	}

	secret := conf.Get(clientSecretKey, "")
	if secret == "" {
		return nil, fmt.Errorf("client secret value is empty, set it: %s", clientSecretKey)
	}

	ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	provider, err := oidc2.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to query provider %q: %v", issuer, err)
	}

	url, err := redirectUrl(conf)
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

func redirectUrl(conf fig.Properties) (string, error) {
	//port := conf.Get("neve.server.port", "")
	//if port == "" {
	//	return "", fmt.Errorf("server port is empty")
	//}
	url := conf.Get("neve.auth.router.callback", "")
	if url == "" {
		return "", fmt.Errorf("callback url is empty")
	}
	//return fmt.Sprintf("https://"), nil

	addr := conf.Get("neve.auth.dex.externalAddr", "")
	if addr == "" {
		return "", fmt.Errorf("external address is empty")
	}
	return fmt.Sprintf("%s/%s", addr, url), nil
}
