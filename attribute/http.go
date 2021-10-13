// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package attribute

import (
	"context"
	"fmt"
	"github.com/xfali/neve-auth/user"
	"net/http"
)

type userKey struct{}
var ctxUserKey userKey

type reqInfoKey struct{}
var ctxReqInfoKeyKey reqInfoKey

type RequestInfo struct {
	Resource string
	Action   string
}

func WithUser(req *http.Request, userInfo *user.UserInfo) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), ctxUserKey, userInfo))
}

func GetUser(req *http.Request) (*user.UserInfo, bool) {
	v := req.Context().Value(ctxUserKey)
	if v != nil {
		return v.(*user.UserInfo), true
	}
	return nil, false
}

func WithRequestInfo(req *http.Request, info *RequestInfo) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), ctxReqInfoKeyKey, info))
}

func GetRequestInfo(req *http.Request) (*RequestInfo, bool) {
	v := req.Context().Value(ctxReqInfoKeyKey)
	if v != nil {
		return v.(*RequestInfo), true
	}
	return nil, false
}

func ParseAttribute(req *http.Request) (Attribute, bool, error) {
	user, ok := GetUser(req)
	if !ok {
		return nil, true, fmt.Errorf("cannot get user info")
	}

	info, ok := GetRequestInfo(req)
	if !ok {
		return nil, true, fmt.Errorf("cannot get request info")
	}

	res := info.Resource
	act := info.Action

	return &DefaultAttribute{
		User:     user,
		Resource: res,
		Action:   act,
	}, true, nil
}
