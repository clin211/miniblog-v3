package errorx

import (
	"fmt"
	"sync"
)

// Range 表示模块错误码的闭区间范围 [Start, End]。
type Range struct {
	Start Code
	End   Code
}

// Registry 维护模块到码段的映射，并提供模块内错误构造能力。
type Registry struct {
	mu      sync.RWMutex
	modules map[string]Range
}

// defaultRegistry 为全局默认注册器。
var defaultRegistry = &Registry{modules: make(map[string]Range)}

// RegisterModule 在全局注册器登记模块的码段范围。
// 若模块已存在或范围非法（Start<=0、End<=Start）则返回错误。
func RegisterModule(module string, r Range) error {
	return defaultRegistry.RegisterModule(module, r)
}

// RegisterModule 为实例方法，登记模块码段。
func (rg *Registry) RegisterModule(module string, r Range) error {
	rg.mu.Lock()
	defer rg.mu.Unlock()

	if module == "" {
		return fmt.Errorf("module 名称不能为空")
	}
	if r.Start <= 0 || r.End <= 0 || r.End <= r.Start {
		return fmt.Errorf("非法的码段范围: [%d, %d]", r.Start, r.End)
	}
	if _, exists := rg.modules[module]; exists {
		return fmt.Errorf("模块已注册: %s", module)
	}
	rg.modules[module] = r
	return nil
}

// NewModuleError 使用全局注册器在模块内构造错误。
// code 必须落在模块登记的码段范围内，否则返回 *Error(CodeBadRequest,... )，提示越界。
func NewModuleError(module string, code Code, message string, opts ...Option) *Error {
	return defaultRegistry.NewModuleError(module, code, message, opts...)
}

// NewModuleError 为实例方法，在模块内构造错误并校验是否在码段范围内。
func (rg *Registry) NewModuleError(module string, code Code, message string, opts ...Option) *Error {
	rg.mu.RLock()
	r, exists := rg.modules[module]
	rg.mu.RUnlock()

	if !exists {
		return New(CodeBadRequest, fmt.Sprintf("未注册的模块: %s", module), opts...)
	}
	if code < r.Start || code > r.End {
		return New(CodeBadRequest, fmt.Sprintf("错误码越界: 模块 %s 范围[%d,%d]，得到 %d", module, r.Start, r.End, code), opts...)
	}
	return New(code, message, opts...)
}
