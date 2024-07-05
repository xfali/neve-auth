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
