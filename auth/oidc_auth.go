// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
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

func (m *OidcLoginMgr) Redirect(ctx *gin.Context) {
	var scopes []string
	if extraScopes := ctx.Query("extra_scopes"); extraScopes != "" {
		scopes = strings.Split(extraScopes, " ")
	}
	var clients []string
	if crossClients := ctx.Query("cross_client"); crossClients != "" {
		clients = strings.Split(crossClients, " ")
	}
	for _, client := range clients {
		scopes = append(scopes, "audience:server:client_id:"+client)
	}
	connectorID := ""
	if id := ctx.Query("connector_id"); id != "" {
		connectorID = id
	}

	authCodeURL := ""
	scopes = append(scopes, "openid", "profile", "email")
	if ctx.Query("offline_access") != "yes" {
		authCodeURL = m.oauth2Config(scopes).AuthCodeURL(ctx.Request.URL.String())
	} else if m.offlineAsScope {
		scopes = append(scopes, "offline_access")
		authCodeURL = m.oauth2Config(scopes).AuthCodeURL(ctx.Request.URL.String())
	} else {
		authCodeURL = m.oauth2Config(scopes).AuthCodeURL(ctx.Request.URL.String(), oauth2.AccessTypeOffline)
	}
	if connectorID != "" {
		authCodeURL = authCodeURL + "&connector_id=" + connectorID
	}

	ctx.Redirect(http.StatusSeeOther, authCodeURL)
}

func (m *OidcLoginMgr) Callback(ctx *gin.Context) {
	var (
		state string
		err   error
		token *oauth2.Token
	)

	// Authorization redirect callback from OAuth2 auth flow.
	if errMsg := ctx.Query("error"); errMsg != "" {
		m.logger.Errorln(errMsg + ": " + ctx.Query("error_description"))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	code := ctx.Query("code")
	if code == "" {
		m.logger.Errorf("no code in request: %q", ctx.Request.RequestURI)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if state = ctx.Query("state"); state == "" {
		m.logger.Errorf("no state in request: %q", ctx.Request.RequestURI)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cctx := oidc.ClientContext(context.Background(), m.client)
	token, err = m.oauth2Config(nil).Exchange(cctx, code)
	if err != nil {
		m.logger.Errorf("failed to get token: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	destToken, err := m.convertToken(token)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = m.writer.WriteToken(ctx.Writer, destToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, state)
}

func (m *OidcLoginMgr) Refresh(ctx *gin.Context) {
	t, err := m.reader.ReadToken(ctx.Request)
	if err != nil {
		m.logger.Errorln(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if t.Refresh == "" {
		m.logger.Errorln("refresh token is empty")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	cctx := oidc.ClientContext(context.Background(), m.client)
	token, err := m.config.TokenSource(cctx, &oauth2.Token{
		RefreshToken: t.Refresh,
	}).Token()
	if err != nil {
		m.logger.Errorf("failed to get token: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	destToken, err := m.convertToken(token)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = m.writer.WriteToken(ctx.Writer, destToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Writer.WriteHeader(http.StatusCreated)
}

func (m *OidcLoginMgr) GetUserInfo(ctx *gin.Context) {

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
