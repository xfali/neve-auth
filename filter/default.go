// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/xfali/fig"
	"github.com/xfali/neve-auth/attribute"
	"github.com/xfali/neve-auth/auth"
	"github.com/xfali/neve-auth/constants"
	"github.com/xfali/xlog"
	"net/http"
)

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
	u, err := constants.RedirectUrl(conf)
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

	attribute.WithUser(req, user)
	return true, nil
}

func (h *defaultHandler) OnFailed(err error, resp http.ResponseWriter, req *http.Request) {
	http.Redirect(resp, req, h.redirectUrl, http.StatusSeeOther)
}

func (h *defaultHandler) ParseAttribute(req *http.Request) (attribute.Attribute, bool, error) {
	return attribute.ParseAttribute(req)
}

type defaultUnauth struct {
}

func NewUnauthHandler() *defaultUnauth {
	return &defaultUnauth{}
}

func (h *defaultUnauth) OnFailed(err error, resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusUnauthorized)
}
