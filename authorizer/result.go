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

package authorizer

const (
	DecisionDeny Decision = iota
	DecisionAllow
	DecisionNoOpinion
)

type Decision int

type Result interface {
	Decision() Decision
	Error() error
}

type defaultResult struct {
	d   Decision
	err error
}

func Allow() *defaultResult {
	return &defaultResult{
		d:   DecisionAllow,
		err: nil,
	}
}

func Deny(err error) *defaultResult {
	return &defaultResult{
		d:   DecisionDeny,
		err: err,
	}
}

func NoOpinion(errs ...error) *defaultResult {
	var err error = nil
	if len(errs) > 0 {
		err = errs[0]
	}
	return &defaultResult{
		d:   DecisionNoOpinion,
		err: err,
	}
}

func (d Decision) IsAllow() bool {
	return d == DecisionAllow
}

func (d Decision) IsDeny() bool {
	return d == DecisionDeny
}

func (r *defaultResult) Decision() Decision {
	return r.d
}

func (r *defaultResult) Error() error {
	return r.err
}
