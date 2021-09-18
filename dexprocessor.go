// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package neveauth

import (
	"github.com/xfali/fig"
	"github.com/xfali/neve-core/bean"
)

type dexOpts struct {
}

type dexProcessor struct {
	dexOpts
}

type DexOpt func(*dexOpts)

func NewDexProcessor(opts ...DexOpt) *dexProcessor {
	ret := &dexProcessor{}
	for _, opt := range opts {
		opt(&ret.dexOpts)
	}
	return ret
}

// 初始化对象处理器
func (p *dexProcessor) Init(conf fig.Properties, container bean.Container) error {

	return nil
}

// 对象分类，判断对象是否实现某些接口，并进行相关归类。为了支持多协程处理，该方法应线程安全。
// 注意：该方法建议只做归类，具体处理使用Process，不保证Processor的实现在此方法中做了相关处理。
// 该方法在Bean Inject注入之后调用
// return: bool 是否能够处理对象， error 处理是否有错误
func (p *dexProcessor) Classify(o interface{}) (bool, error) {
	return false, nil
}

// 对已分类对象做统一处理，注意如果存在耗时操作，请使用其他协程处理。
// 该方法在Classify及BeanAfterSet后调用。
// 成功返回nil，失败返回error
func (p *dexProcessor) Process() error {
	return nil
}

func (p *dexProcessor) BeanDestroy() error {
	return nil
}
