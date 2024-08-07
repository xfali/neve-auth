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

package oidc

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/xfali/neve-auth/token"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type OidcContext struct {
	IssuerURL string

	config         *oauth2.Config
	provider       *oidc.Provider
	verifier       *oidc.IDTokenVerifier
	client         *http.Client
	offlineAsScope bool
}

func NewOidcContext(client *http.Client, issuerURL string, config *oauth2.Config) *OidcContext {
	if client == nil {
		client = http.DefaultClient
	}
	ret := &OidcContext{
		IssuerURL: issuerURL,
		config:    config,
		client:    client,
	}

	return ret
}

func (c *OidcContext) Init() error {
	ctx := oidc.ClientContext(context.Background(), c.client)
	provider, err := oidc.NewProvider(ctx, c.IssuerURL)
	if err != nil {
		return fmt.Errorf("failed to query provider %q: %v", c.IssuerURL, err)
	}
	c.provider = provider

	var s struct {
		// What scopes does a provider support?
		//
		// See: https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
		ScopesSupported []string `json:"scopes_supported"`
	}
	if err := provider.Claims(&s); err != nil {
		return fmt.Errorf("failed to parse provider scopes_supported: %v", err)
	}

	if len(s.ScopesSupported) == 0 {
		// scopes_supported is a "RECOMMENDED" discovery claim, not a required
		// one. If missing, assume that the provider follows the spec and has
		// an "offline_access" scope.
		c.offlineAsScope = true
	} else {
		// See if scopes_supported has the "offline_access" scope.
		c.offlineAsScope = func() bool {
			for _, scope := range s.ScopesSupported {
				if scope == oidc.ScopeOfflineAccess {
					return true
				}
			}
			return false
		}()
	}
	return nil
}

func (c *OidcContext) Oauth2Config(scopes []string) *oauth2.Config {
	if len(scopes) == 0 {
		return c.config
	}
	ret := *c.config
	ret.Scopes = scopes
	return &ret
}

func (c *OidcContext) GetProviderConfig(ctx context.Context) error {
	wellKnown := strings.TrimSuffix(c.IssuerURL, "/") + "/.well-known/openid-configuration"
	_, err := http.NewRequest("GET", wellKnown, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *OidcContext) GetProvider(ctx context.Context) *oidc.Provider {
	return c.provider
}

func (c *OidcContext) GetVerifier(ctx context.Context) token.Verifier {
	v := c.provider.Verifier(&oidc.Config{ClientID: c.config.ClientID})
	return token.NewOidcVerifier(v)
}

//func (c *OidcContext) SetupClient(store st) error{
//
//}
