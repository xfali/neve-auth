// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/xfali/neve-auth/attribute"
	"github.com/xfali/neve-auth/auth"
	"github.com/xfali/xlog"
	"net/http"
)

type defaultHandler struct {
	logger          xlog.Logger
	matcher         PathMatcher
	reader          auth.TokenReader
	redirectHandler auth.RedirectHandler
}

func NewDefaultHandler(matcher PathMatcher, redirectHandler auth.RedirectHandler, reader auth.TokenReader) *defaultHandler {
	ret := &defaultHandler{
		logger:          xlog.GetLogger(),
		matcher:         matcher,
		reader:          reader,
		redirectHandler: redirectHandler,
	}

	return ret
}

func (h *defaultHandler) ParseToken(req *http.Request) (string, bool, error) {
	// ignore
	if !h.matcher.Match(req.RequestURI) {
		return "", false, nil
	}
	token, err := h.reader.ReadToken(req)
	if err != nil {
		return "", false, err
	}
	return token.ID, true, nil
}

func (h *defaultHandler) OnFailed(err error, resp http.ResponseWriter, req *http.Request) {
	h.redirectHandler.Redirect(resp, req)
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
