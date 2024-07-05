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

package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/xfali/xlog"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"time"
)

var (
	cookieName = "neve-token"
)

type OidcLoginMgr struct {
	logger         xlog.Logger
	config         *oauth2.Config
	offlineAsScope bool
	client         *http.Client

	reader TokenReader
	writer TokenWriter
}

func NewOidcLoginMgr(config *oauth2.Config, client *http.Client) *OidcLoginMgr {
	ret := &OidcLoginMgr{
		logger: xlog.GetLogger(),
		config: config,
		reader: NewTokenReader(),
		writer: NewTokenWriter(),
		client: client,
	}
	return ret
}

func (m *OidcLoginMgr) Redirect(w http.ResponseWriter, r *http.Request) {
	var scopes []string
	if extraScopes := r.FormValue("extra_scopes"); extraScopes != "" {
		scopes = strings.Split(extraScopes, " ")
	}
	var clients []string
	if crossClients := r.FormValue("cross_client"); crossClients != "" {
		clients = strings.Split(crossClients, " ")
	}
	for _, client := range clients {
		scopes = append(scopes, "audience:server:client_id:"+client)
	}
	connectorID := ""
	if id := r.FormValue("connector_id"); id != "" {
		connectorID = id
	}

	authCodeURL := ""
	scopes = append(scopes, "openid", "profile", "email")
	if r.FormValue("offline_access") != "yes" {
		authCodeURL = m.oauth2Config(scopes).AuthCodeURL(r.URL.String())
	} else if m.offlineAsScope {
		scopes = append(scopes, "offline_access")
		authCodeURL = m.oauth2Config(scopes).AuthCodeURL(r.URL.String())
	} else {
		authCodeURL = m.oauth2Config(scopes).AuthCodeURL(r.URL.String(), oauth2.AccessTypeOffline)
	}
	if connectorID != "" {
		authCodeURL = authCodeURL + "&connector_id=" + connectorID
	}

	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
}

func (m *OidcLoginMgr) Callback(w http.ResponseWriter, r *http.Request) {
	var (
		state string
		err   error
		token *oauth2.Token
	)

	// Authorization redirect callback from OAuth2 auth flow.
	if errMsg := r.FormValue("error"); errMsg != "" {
		m.logger.Errorln(errMsg + ": " + r.FormValue("error_description"))
		http.Error(w, errMsg+": "+r.FormValue("error_description"), http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	if code == "" {
		m.logger.Errorf("no code in request: %q", r.RequestURI)
		http.Error(w, fmt.Sprintf("no code in request: %q", r.Form), http.StatusBadRequest)
		return
	}
	if state = r.FormValue("state"); state == "" {
		m.logger.Errorf("no state in request: %q", r.RequestURI)
		http.Error(w, fmt.Sprintf("expected state %q", r.RequestURI), http.StatusBadRequest)
		return
	}

	cctx := oidc.ClientContext(context.Background(), m.client)
	token, err = m.oauth2Config(nil).Exchange(cctx, code)
	if err != nil {
		m.logger.Errorf("failed to get token: %v", err)
		http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
		return
	}

	destToken, err := m.convertToken(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to convert token: %v", err), http.StatusInternalServerError)
		return
	}
	err = m.writer.WriteToken(w, destToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to response write token: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, state, http.StatusFound)
}

func (m *OidcLoginMgr) Refresh(w http.ResponseWriter, r *http.Request) {
	t, err := m.reader.ReadToken(r)
	if err != nil {
		m.logger.Errorln(err)
		http.Error(w, fmt.Sprintf("no refresh_token in request: %v", err), http.StatusBadRequest)
		return
	}

	if t.Refresh == "" {
		m.logger.Errorln("refresh token is empty")
		http.Error(w, fmt.Sprintf("no refresh_token in request"), http.StatusBadRequest)
		return
	}

	cctx := oidc.ClientContext(context.Background(), m.client)
	token, err := m.config.TokenSource(cctx, &oauth2.Token{
		RefreshToken: t.Refresh,
	}).Token()
	if err != nil {
		m.logger.Errorf("failed to get token: %v", err)
		http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
		return
	}

	destToken, err := m.convertToken(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to convert token: %v", err), http.StatusInternalServerError)
		return
	}

	err = m.writer.WriteToken(w, destToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to response write token: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (m *OidcLoginMgr) GetUserInfo(w http.ResponseWriter, r *http.Request) {

}

func (m *OidcLoginMgr) oauth2Config(scopes []string) *oauth2.Config {
	if len(scopes) == 0 {
		return m.config
	}
	ret := *m.config
	ret.Scopes = scopes
	return &ret
}

func (m *OidcLoginMgr) convertToken(token *oauth2.Token) (*Token, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		err := fmt.Errorf("no id_token in token response")
		m.logger.Errorln(err)
		return nil, err
	}

	return &Token{
		ID:      rawIDToken,
		Refresh: token.RefreshToken,
		Expire:  token.Expiry,
	}, nil
}

type defaultTokenWriter struct{}

func NewTokenWriter() *defaultTokenWriter {
	return &defaultTokenWriter{}
}

func (m *defaultTokenWriter) WriteToken(resp http.ResponseWriter, token *Token) error {
	tokenData, err := json.Marshal(token)
	if err != nil {
		return err
	}
	tokenStr := base64.StdEncoding.EncodeToString(tokenData)
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    tokenStr,
		MaxAge:   int(time.Until(token.Expire).Seconds()),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(resp, cookie)
	return nil
}

type defaultTokenReader struct{}

func NewTokenReader() *defaultTokenReader {
	return &defaultTokenReader{}
}

func (m *defaultTokenReader) ReadToken(req *http.Request) (*Token, error) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return nil, fmt.Errorf("get cookie failed: %v", err)
	}

	tokenData, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("decode token failed: %v", err)
	}
	var token Token
	err = json.Unmarshal(tokenData, &token)
	if err != nil {
		return nil, fmt.Errorf("decode token failed: %v", err)
	}

	if token.ID == "" || token.Expire.Before(time.Now()) {
		err := fmt.Errorf("token invalid. id: %s time: %v", token.ID, token.Expire)
		return nil, err
	}
	return &token, nil
}
