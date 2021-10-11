// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package config

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/casbin/gorm-adapter/v2"
	"github.com/casbin/redis-adapter/v2"
	"github.com/xfali/fig"
	"github.com/xfali/neve-core/bean"
	"strings"
)

const (
	casbinModelKey        = "neve.auth.casbin.modelFile"
	casbinAdapterTypeKey  = "neve.auth.casbin.adapter.type"
	casbinAdapterValueKey = "neve.auth.casbin.adapter.value"
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
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.obj

[policy_effect]
e = some(where (p.eft == allow))
`
)

type CasbinConfig struct {}

func NewCasbinConfig() *CasbinConfig {
	return &CasbinConfig{}
}

func (c *CasbinConfig) Init(conf fig.Properties, container bean.Container) (*casbin.Enforcer, error) {
	var m model.Model
	var err error
	modelFile := conf.Get(casbinModelKey, "")
	if modelFile != "" {
		m, err = model.NewModelFromFile(modelFile)
	} else {
		m, err = model.NewModelFromString(defaultModel)
	}

	if err != nil {
		return nil, err
	}

	t := conf.Get(casbinAdapterTypeKey, "")
	v := conf.Get(casbinAdapterValueKey, "")
	if t == "" || v == "" {
		return nil, fmt.Errorf("casbin adapter type: %s value: %s . ", t, v)
	}
	adapter, err := selectAdapter(t, v)
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}

func selectAdapter(t, v string) (persist.Adapter, error) {
	switch t {
	case "file":
		return fileadapter.NewAdapter(v), nil
	case "mysql":
		return gormadapter.NewAdapter("mysql", v)
	case "posgresql":
		return gormadapter.NewAdapter("posgresql", v)
	case "redis":
		return redisadapter.NewAdapter("tcp", v), nil
	default:
		if len(t) > 3 && strings.ToLower(t[:3]) == "db:"{
			return gormadapter.NewAdapter(t[3:], v)
		}
	}
	return nil, nil
}
