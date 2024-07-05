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
		logger:   xlog.GetLogger(),
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
