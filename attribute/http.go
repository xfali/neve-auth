// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package attribute

import (
	"context"
	"fmt"
	"github.com/xfali/neve-auth/constants"
	"github.com/xfali/neve-auth/user"
	"net/http"
)

type userKey struct{}

var reqUserKey userKey

func WithUser(req *http.Request, userInfo *user.UserInfo) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), reqUserKey, userInfo))
}

func GetUser(req *http.Request) (*user.UserInfo, bool) {
	v := req.Context().Value(reqUserKey)
	if v != nil {
		return v.(*user.UserInfo), true
	}
	return nil, false
}

func GetResource(req *http.Request) string {
	return req.Header.Get(constants.HeaderKeyResource)
}

func GetAction(req *http.Request) string {
	return req.Header.Get(constants.HeaderKeyAction)
}

func ParseAttribute(req *http.Request) (Attribute, bool, error) {
	user, ok := GetUser(req)
	if !ok {
		return nil, true, fmt.Errorf("cannot get user info")
	}

	res := GetResource(req)
	act := GetAction(req)

	return &DefaultAttribute{
		User:     user,
		Resource: res,
		Action:   act,
	}, true, nil
}
