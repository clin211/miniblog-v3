// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package validate

import (
	"fmt"
	"strings"
)

// ErrorCode 错误代码类型
type ErrorCode string

// 预定义错误代码
const (
	ErrorCodeRequired         ErrorCode = "REQUIRED"
	ErrorCodeInvalidEmail     ErrorCode = "INVALID_EMAIL"
	ErrorCodeInvalidURL       ErrorCode = "INVALID_URL"
	ErrorCodeInvalidUUID      ErrorCode = "INVALID_UUID"
	ErrorCodeInvalidJSON      ErrorCode = "INVALID_JSON"
	ErrorCodeInvalidAlpha     ErrorCode = "INVALID_ALPHA"
	ErrorCodeInvalidNumeric   ErrorCode = "INVALID_NUMERIC"
	ErrorCodeInvalidAlphanum  ErrorCode = "INVALID_ALPHANUM"
	ErrorCodeInvalidASCII     ErrorCode = "INVALID_ASCII"
	ErrorCodeInvalidIPv4      ErrorCode = "INVALID_IPV4"
	ErrorCodeInvalidIPv6      ErrorCode = "INVALID_IPV6"
	ErrorCodeInvalidPort      ErrorCode = "INVALID_PORT"
	ErrorCodeInvalidDNS       ErrorCode = "INVALID_DNS"
	ErrorCodeInvalidHost      ErrorCode = "INVALID_HOST"
	ErrorCodeInvalidMAC       ErrorCode = "INVALID_MAC"
	ErrorCodeInvalidBase64    ErrorCode = "INVALID_BASE64"
	ErrorCodeInvalidDataURI   ErrorCode = "INVALID_DATA_URI"
	ErrorCodeInvalidRFC3339   ErrorCode = "INVALID_RFC3339"
	ErrorCodeInvalidSemver    ErrorCode = "INVALID_SEMVER"
	ErrorCodeInvalidULID      ErrorCode = "INVALID_ULID"
	ErrorCodeInvalidYYYYMMDD  ErrorCode = "INVALID_YYYYMMDD"
	ErrorCodeInvalidLatitude  ErrorCode = "INVALID_LATITUDE"
	ErrorCodeInvalidLongitude ErrorCode = "INVALID_LONGITUDE"
	ErrorCodeInvalidSSN       ErrorCode = "INVALID_SSN"
	ErrorCodeInvalidLength    ErrorCode = "INVALID_LENGTH"
	ErrorCodeInvalidRange     ErrorCode = "INVALID_RANGE"
	ErrorCodeInvalidMatches   ErrorCode = "INVALID_MATCHES"
	ErrorCodeInvalidEnum      ErrorCode = "INVALID_ENUM"
	ErrorCodeUnknownRule      ErrorCode = "UNKNOWN_RULE"
)

// ValidationErrorWithCode 带错误代码的验证错误
type ValidationErrorWithCode struct {
	Code    ErrorCode   `json:"code"`
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
	Rule    string      `json:"rule"`
}

// Error 实现 error 接口
func (e *ValidationErrorWithCode) Error() string {
	return fmt.Sprintf("字段 '%s' 验证失败 [%s]: %s (值: %v, 规则: %s)",
		e.Field, e.Code, e.Message, e.Value, e.Rule)
}

// ValidationErrorsWithCode 多个带错误代码的验证错误
type ValidationErrorsWithCode []*ValidationErrorWithCode

// Error 实现 error 接口
func (es ValidationErrorsWithCode) Error() string {
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
func (es ValidationErrorsWithCode) HasErrors() bool {
	return len(es) > 0
}

// GetErrorsByCode 根据错误代码获取错误列表
func (es ValidationErrorsWithCode) GetErrorsByCode(code ErrorCode) ValidationErrorsWithCode {
	var result ValidationErrorsWithCode
	for _, err := range es {
		if err.Code == code {
			result = append(result, err)
		}
	}
	return result
}

// GetErrorsByField 根据字段名获取错误列表
func (es ValidationErrorsWithCode) GetErrorsByField(field string) ValidationErrorsWithCode {
	var result ValidationErrorsWithCode
	for _, err := range es {
		if err.Field == field {
			result = append(result, err)
		}
	}
	return result
}

// ToMap 转换为错误映射
func (es ValidationErrorsWithCode) ToMap() map[string][]string {
	result := make(map[string][]string)
	for _, err := range es {
		result[err.Field] = append(result[err.Field], err.Message)
	}
	return result
}

// NewValidationError 创建新的验证错误
func NewValidationError(code ErrorCode, field, message string, value interface{}, rule string) *ValidationErrorWithCode {
	return &ValidationErrorWithCode{
		Code:    code,
		Field:   field,
		Message: message,
		Value:   value,
		Rule:    rule,
	}
}

// NewRequiredError 创建必填字段错误
func NewRequiredError(field string, value interface{}) *ValidationErrorWithCode {
	return NewValidationError(
		ErrorCodeRequired,
		field,
		"字段不能为空",
		value,
		"required",
	)
}

// NewEmailError 创建邮箱格式错误
func NewEmailError(field, value string) *ValidationErrorWithCode {
	return NewValidationError(
		ErrorCodeInvalidEmail,
		field,
		"邮箱格式不正确",
		value,
		"email",
	)
}

// NewLengthError 创建长度错误
func NewLengthError(field, value string, min, max int) *ValidationErrorWithCode {
	var message string
	if min == max {
		message = fmt.Sprintf("长度必须为 %d 个字符", min)
	} else {
		message = fmt.Sprintf("长度必须在 %d 到 %d 个字符之间", min, max)
	}

	return NewValidationError(
		ErrorCodeInvalidLength,
		field,
		message,
		value,
		fmt.Sprintf("length(%d|%d)", min, max),
	)
}

// NewRangeError 创建范围错误
func NewRangeError(field string, value interface{}, min, max float64) *ValidationErrorWithCode {
	return NewValidationError(
		ErrorCodeInvalidRange,
		field,
		fmt.Sprintf("值必须在 %v 到 %v 之间", min, max),
		value,
		fmt.Sprintf("range(%v|%v)", min, max),
	)
}

// NewMatchesError 创建正则匹配错误
func NewMatchesError(field, value, pattern string) *ValidationErrorWithCode {
	return NewValidationError(
		ErrorCodeInvalidMatches,
		field,
		fmt.Sprintf("值不符合正则表达式模式: %s", pattern),
		value,
		fmt.Sprintf("matches(%s)", pattern),
	)
}

// NewEnumError 创建枚举值错误
func NewEnumError(field string, value interface{}, allowedValues []string) *ValidationErrorWithCode {
	return NewValidationError(
		ErrorCodeInvalidEnum,
		field,
		fmt.Sprintf("值必须是以下之一: %s", strings.Join(allowedValues, ", ")),
		value,
		fmt.Sprintf("in(%s)", strings.Join(allowedValues, "|")),
	)
}

// NewUnknownRuleError 创建未知规则错误
func NewUnknownRuleError(field, rule string, value interface{}) *ValidationErrorWithCode {
	return NewValidationError(
		ErrorCodeUnknownRule,
		field,
		fmt.Sprintf("未知的验证规则: %s", rule),
		value,
		rule,
	)
}
