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
