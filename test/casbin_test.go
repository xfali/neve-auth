// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

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
	ok ,err := enforcer.Enforce(subject, project, resource, action)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expect true get : ", ok)
	}
}
