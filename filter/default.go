// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"context"
	"fmt"
	"github.com/xfali/fig"
	"github.com/xfali/neve-auth/auth"
	"github.com/xfali/neve-auth/authorizer"
	"github.com/xfali/neve-auth/config"
	"github.com/xfali/neve-auth/user"
	"github.com/xfali/xlog"
	"net/http"
)

type userKey struct{}

var reqUserKey userKey

func WithUser(req *http.Request, userInfo *user.UserInfo) {
	req.WithContext(context.WithValue(req.Context(), reqUserKey, userInfo))
}

func GetUser(req *http.Request) (*user.UserInfo, bool) {
	v := req.Context().Value(reqUserKey)
	if v != nil {
		return v.(*user.UserInfo), true
	}
	return nil, false
}

type defaultHandler struct {
	logger        xlog.Logger
	reader        auth.TokenReader
	authenticator Authenticator
	redirectUrl   string
}

func NewDefaultHandler(conf fig.Properties, reader auth.TokenReader, authenticator Authenticator) *defaultHandler {
	ret := &defaultHandler{
		logger:        xlog.GetLogger(),
		reader:        reader,
		authenticator: authenticator,
	}
	u, err := config.RedirectUrl(conf)
	if err != nil {
		ret.logger.Errorln(err)
		return nil
	}
	ret.redirectUrl = u

	return ret
}

func (h *defaultHandler) ParseToken(req *http.Request) (bool, error) {
	token, err := h.reader.ReadToken(req)
	if err != nil {
		return false, err
	}
	user, err := h.authenticator.AuthenticateToken(req.Context(), token.ID)
	if err != nil {
		return false, err
	}

	WithUser(req, user)
	return true, nil
}

func (h *defaultHandler) OnFailed(err error, resp http.ResponseWriter, req *http.Request) {
	http.Redirect(resp, req, h.redirectUrl, http.StatusSeeOther)
}

func (h *defaultHandler) ParseAttribute(req *http.Request) (authorizer.Attribute, bool, error) {
	user, ok := GetUser(req)
	if !ok {
		return nil, true, fmt.Errorf("cannot get user info")
	}


}

type defaultUnauth struct {
}

func (h *defaultUnauth) OnFailed(err error, resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusUnauthorized)
}