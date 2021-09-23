// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

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

func (r *defaultResult) Decision() Decision {
	return r.d
}

func (r *defaultResult) Error() error {
	return r.err
}
