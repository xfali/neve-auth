// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package neveauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dexidp/dex/storage"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type defaultAuthenticator struct {
	clientId       string
	issuerURL      string
	provider       *oidc.Provider
	verifier       *oidc.IDTokenVerifier
	scopes         []string
	client         *http.Client
	offlineAsScope bool
}

func CreateAuthenticator(client *http.Client, clientId string, issuerURL string) (*defaultAuthenticator, error) {
	ret := &defaultAuthenticator{
		client:    client,
		issuerURL: issuerURL,
	}
	ctx := oidc.ClientContext(context.Background(), client)
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query provider %q: %v", issuerURL, err)
	}

	var s struct {
		// What scopes does a provider support?
		//
		// See: https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
		ScopesSupported []string `json:"scopes_supported"`
	}
	if err := provider.Claims(&s); err != nil {
		return nil, fmt.Errorf("failed to parse provider scopes_supported: %v", err)
	}

	if len(s.ScopesSupported) == 0 {
		// scopes_supported is a "RECOMMENDED" discovery claim, not a required
		// one. If missing, assume that the provider follows the spec and has
		// an "offline_access" scope.
		ret.offlineAsScope = true
	} else {
		// See if scopes_supported has the "offline_access" scope.
		ret.offlineAsScope = func() bool {
			for _, scope := range s.ScopesSupported {
				if scope == oidc.ScopeOfflineAccess {
					return true
				}
			}
			return false
		}()
	}

	ret.provider = provider
	ret.verifier = provider.Verifier(&oidc.Config{ClientID: clientId})

	return ret, nil
}

func (a *defaultAuthenticator) AuthenticateToken(ctx context.Context, token string) error {
	err := a.VerifyIssuer(token)
	if err != nil {
		return err
	}

	idToken, err := a.verifier.Verify(ctx, token)
	if err != nil {
		return tokenVerifyError.V(err)
	}

	var c storage.Claims
	if err := idToken.Claims(&c); err != nil {
		return parseClaimsError.V(err)
	}

}

func (a *defaultAuthenticator) VerifyIssuer(token string) error {
	issuer, err := parseIssuer(token)
	if err != nil {
		return err
	}
	if issuer != a.issuerURL {
		return issuerVerifyError
	}
	return nil
}

func (a *defaultAuthenticator) oauth2Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     a.clientID,
		ClientSecret: a.clientSecret,
		Endpoint:     a.provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  a.redirectURI,
	}
}

func parseIssuer(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", malformedTokenError
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", decodeTokenError.V(err)
	}
	return string(payload), nil
}
