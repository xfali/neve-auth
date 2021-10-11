// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package auth

import (
	"net/http"
)

type TokenReader interface {
	ReadToken(req *http.Request) (*Token, error)
}

type TokenWriter interface {
	WriteToken(resp http.ResponseWriter, token *Token) error
}

type RedirectHandler interface {
	Redirect(resp http.ResponseWriter, req *http.Request)
}

type CallbackHandler interface {
	Callback(resp http.ResponseWriter, req *http.Request)
}

type RefreshHandler interface {
	Refresh(resp http.ResponseWriter, req *http.Request)
}

type UserInfoHandler interface {
	GetUserInfo(resp http.ResponseWriter, req *http.Request)
}
