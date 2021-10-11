// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/xfali/fig"
	"github.com/xfali/neve-auth/constants"
	"github.com/xfali/router"
	"github.com/xfali/xlog"
)

type PathMatcher interface {
	Include(paths ...string) PathMatcher
	Exclude(paths ...string) PathMatcher

	Match(path string) bool
}

type DefaultPathMatcher struct {
	logger xlog.Logger

	includes router.Router
	excludes router.Router
}

func NewPathMatcher(config fig.Properties) *DefaultPathMatcher {
	ret := &DefaultPathMatcher{
		logger: xlog.GetLogger(),
		includes: router.New(),
		excludes: router.New(),
	}

	var includes, excludes []string
	err := config.GetValue(constants.IncludesKey, &includes)
	if err != nil {
		ret.logger.Errorln(err)
	}
	ret.Include(includes...)

	err = config.GetValue(constants.ExcludesKey, &excludes)
	if err != nil {
		ret.logger.Errorln(err)
	}
	ret.Exclude(excludes...)
	return ret
}

func (f *DefaultPathMatcher) Include(paths ...string) PathMatcher {
	for _, v := range paths {
		_ = f.includes.AddRoute(v, nil)
	}
	return f
}

func (f *DefaultPathMatcher) Exclude(paths ...string) PathMatcher {
	for _, v := range paths {
		_ = f.excludes.AddRoute(v, nil)
	}
	return f
}

func (f *DefaultPathMatcher) Match(path string) bool {
	if _, err := f.excludes.Match(path, nil); err == nil {
		return false
	}
	if _, err := f.includes.Match(path, nil); err == nil {
		return true
	}
	return false
}
