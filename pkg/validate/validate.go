package validate

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Validator 验证器接口
type Validator interface {
	Validate() error
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

// Error 实现 error 接口
func (e *ValidationError) Error() string {
	return fmt.Sprintf("字段 '%s' 验证失败: %s (值: %v)", e.Field, e.Message, e.Value)
}

// ValidationErrors 多个验证错误
type ValidationErrors []*ValidationError

// Error 实现 error 接口
func (es ValidationErrors) Error() string {
	if len(es) == 0 {
		return ""
	}

	var messages []string
	for _, err := range es {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors 检查是否有错误
func (es ValidationErrors) HasErrors() bool {
	return len(es) > 0
}

// ValidateStruct 验证结构体
func ValidateStruct(v interface{}) error {
	// 使用 govalidator 进行基础验证
	valid, err := govalidator.ValidateStruct(v)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("结构体验证失败")
	}

	return nil
}

// ValidateStructWithCustomRules 使用自定义规则验证结构体
func ValidateStructWithCustomRules(v interface{}) error {
	var errors ValidationErrors

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("只能验证结构体类型")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 获取验证标签
		tag := fieldType.Tag.Get("valid")
		if tag == "" || tag == "-" {
			continue
		}

		// 执行自定义验证
		if err := validateField(field, fieldType, tag); err != nil {
			errors = append(errors, &ValidationError{
				Field:   fieldType.Name,
				Message: err.Error(),
				Value:   field.Interface(),
			})
		}
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// validateField 验证单个字段
func validateField(field reflect.Value, fieldType reflect.StructField, tag string) error {
	// 解析验证标签
	rules := parseValidationRules(tag)

	for _, rule := range rules {
		if err := applyValidationRule(field, fieldType, rule); err != nil {
			return err
		}
	}

	return nil
}

// parseValidationRules 解析验证规则
func parseValidationRules(tag string) []string {
	return strings.Split(tag, ",")
}

// applyValidationRule 应用验证规则
func applyValidationRule(field reflect.Value, fieldType reflect.StructField, rule string) error {
	rule = strings.TrimSpace(rule)

	// 检查是否是带参数的规则
	if strings.Contains(rule, "(") && strings.Contains(rule, ")") {
		return applyParametricRule(field, fieldType, rule)
	}

	// 应用无参数规则
	return applySimpleRule(field, fieldType, rule)
}

// applySimpleRule 应用简单规则
func applySimpleRule(field reflect.Value, fieldType reflect.StructField, rule string) error {
	switch rule {
	case "required":
		return validateRequired(field, fieldType)
	case "email":
		return validateEmail(field, fieldType)
	case "alphanum":
		return validateAlphanum(field, fieldType)
	case "numeric":
		return validateNumeric(field, fieldType)
	case "alpha":
		return validateAlpha(field, fieldType)
	case "ascii":
		return validateASCII(field, fieldType)
	case "url":
		return validateURL(field, fieldType)
	case "ipv4":
		return validateIPv4(field, fieldType)
	case "ipv6":
		return validateIPv6(field, fieldType)
	case "uuid":
		return validateUUID(field, fieldType)
	case "json":
		return validateJSON(field, fieldType)
	default:
		// 尝试使用 govalidator 的内置验证器
		return validateWithGovalidator(field, fieldType, rule)
	}
}

// applyParametricRule 应用带参数的规则
func applyParametricRule(field reflect.Value, fieldType reflect.StructField, rule string) error {
	// 解析规则名称和参数
	ruleName, params := parseParametricRule(rule)

	switch ruleName {
	case "length":
		return validateLength(field, fieldType, params)
	case "range":
		return validateRange(field, fieldType, params)
	case "matches":
		return validateMatches(field, fieldType, params)
	case "in":
		return validateIn(field, fieldType, params)
	default:
		return fmt.Errorf("未知的带参数验证规则: %s", ruleName)
	}
}

// parseParametricRule 解析带参数的规则
func parseParametricRule(rule string) (string, []string) {
	start := strings.Index(rule, "(")
	end := strings.LastIndex(rule, ")")

	if start == -1 || end == -1 || start >= end {
		return rule, nil
	}

	ruleName := strings.TrimSpace(rule[:start])
	paramsStr := rule[start+1 : end]
	params := strings.Split(paramsStr, "|")

	// 清理参数
	for i, param := range params {
		params[i] = strings.TrimSpace(param)
	}

	return ruleName, params
}

// httpxValidatorAdapter 实现 httpx.Validator，用于与 go-zero httpx 适配。
type HttpxValidatorAdapter struct{}

// Validate 实现 httpx.Validator 接口。
func (HttpxValidatorAdapter) Validate(r *http.Request, v interface{}) error {
	return ValidateStructWithCustomRules(v)
}
