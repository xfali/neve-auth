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

package test

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist/file-adapter"
	"testing"
)

const (
	defaultModel = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _,_,_

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act

[policy_effect]
e = some(where (p.eft == allow))
`
)

func TestFileAdapter(t *testing.T) {
	adapter := fileadapter.NewAdapter("assets/rbac_policy.csv")
	m, err := model.NewModelFromString(defaultModel)
	if err != nil {
		t.Fatal(err)
	}
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		t.Fatal(err)
	}

	var subject, project, resource, action = "admin", "tenant1", "userinfo", "read"
	myRes := enforcer.GetPermissionsForUser("admin")
	t.Log("Permissions for ", subject, ": ", myRes)
	ok, err := enforcer.Enforce(subject, project, resource, action)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expect true get : ", ok)
	}
}
