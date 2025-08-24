// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package validate

// 常用验证规则常量
const (
	// 基础验证规则
	RuleRequired = "required"
	RuleEmail    = "email"
	RuleURL      = "url"
	RuleUUID     = "uuid"
	RuleJSON     = "json"

	// 字符类型验证规则
	RuleAlpha    = "alpha"
	RuleNumeric  = "numeric"
	RuleAlphanum = "alphanum"
	RuleASCII    = "ascii"

	// 网络相关验证规则
	RuleIPv4 = "ipv4"
	RuleIPv6 = "ipv6"
	RulePort = "port"
	RuleDNS  = "dns"
	RuleHost = "host"
	RuleMAC  = "mac"

	// 数据格式验证规则
	RuleBase64   = "base64"
	RuleDataURI  = "datauri"
	RuleRFC3339  = "rfc3339"
	RuleSemver   = "semver"
	RuleULID     = "ulid"
	RuleYYYYMMDD = "yyyymmdd"

	// 地理位置验证规则
	RuleLatitude  = "latitude"
	RuleLongitude = "longitude"

	// 其他验证规则
	RuleSSN = "ssn"
)

// 常用长度规则
const (
	// 用户名长度规则
	RuleUsernameLength = "length(3|20)"

	// 密码长度规则
	RulePasswordLength = "length(6|50)"

	// 邮箱长度规则
	RuleEmailLength = "length(5|254)"

	// 手机号长度规则
	RulePhoneLength = "length(11|11)"

	// 验证码长度规则
	RuleCodeLength = "length(4|6)"
)

// 常用数值范围规则
const (
	// 年龄范围规则
	RuleAgeRange = "range(1|120)"

	// 评分范围规则
	RuleScoreRange = "range(0|100)"

	// 价格范围规则
	RulePriceRange = "range(0|999999)"
)

// 常用正则表达式规则
const (
	// 手机号格式规则
	RulePhonePattern = "matches(^1[3-9]\\d{9}$)"

	// 身份证号格式规则
	RuleIDCardPattern = "matches(^[1-9]\\d{5}(18|19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$)"

	// 中文姓名格式规则
	RuleChineseNamePattern = "matches(^[\\u4e00-\\u9fa5]{2,4}$)"

	// 英文姓名格式规则
	RuleEnglishNamePattern = "matches(^[a-zA-Z\\s]{2,50}$)"
)

// 常用枚举值规则
const (
	// 性别枚举规则
	RuleGenderEnum = "in(male|female|other)"

	// 状态枚举规则
	RuleStatusEnum = "in(active|inactive|pending|deleted)"

	// 类型枚举规则
	RuleTypeEnum = "in(user|admin|moderator)"
)

// ValidationRule 验证规则结构体
type ValidationRule struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Example     string   `json:"example"`
	Tags        []string `json:"tags"`
}

// GetCommonRules 获取常用验证规则
func GetCommonRules() map[string]ValidationRule {
	return map[string]ValidationRule{
		// 基础验证规则
		"required": {
			Name:        "required",
			Description: "字段不能为空",
			Example:     `valid:"required"`,
			Tags:        []string{"基础", "必填"},
		},
		"email": {
			Name:        "email",
			Description: "邮箱格式验证",
			Example:     `valid:"required,email"`,
			Tags:        []string{"基础", "邮箱"},
		},
		"url": {
			Name:        "url",
			Description: "URL格式验证",
			Example:     `valid:"url"`,
			Tags:        []string{"基础", "URL"},
		},

		// 字符类型验证规则
		"alpha": {
			Name:        "alpha",
			Description: "只能包含字母",
			Example:     `valid:"alpha"`,
			Tags:        []string{"字符", "字母"},
		},
		"numeric": {
			Name:        "numeric",
			Description: "只能包含数字",
			Example:     `valid:"numeric"`,
			Tags:        []string{"字符", "数字"},
		},
		"alphanum": {
			Name:        "alphanum",
			Description: "只能包含字母和数字",
			Example:     `valid:"alphanum"`,
			Tags:        []string{"字符", "字母数字"},
		},

		// 长度验证规则
		"length": {
			Name:        "length",
			Description: "字符串长度验证",
			Example:     `valid:"length(3|20)"`,
			Tags:        []string{"长度", "范围"},
		},

		// 数值范围验证规则
		"range": {
			Name:        "range",
			Description: "数值范围验证",
			Example:     `valid:"range(1|100)"`,
			Tags:        []string{"数值", "范围"},
		},

		// 正则表达式验证规则
		"matches": {
			Name:        "matches",
			Description: "正则表达式匹配验证",
			Example:     `valid:"matches(^1[3-9]\\d{9}$)"`,
			Tags:        []string{"正则", "模式"},
		},

		// 枚举值验证规则
		"in": {
			Name:        "in",
			Description: "枚举值验证",
			Example:     `valid:"in(active|inactive|pending)"`,
			Tags:        []string{"枚举", "选择"},
		},
	}
}

// GetRuleByName 根据规则名称获取规则信息
func GetRuleByName(name string) (ValidationRule, bool) {
	rules := GetCommonRules()
	rule, exists := rules[name]
	return rule, exists
}

// GetRulesByTag 根据标签获取规则列表
func GetRulesByTag(tag string) []ValidationRule {
	rules := GetCommonRules()
	var result []ValidationRule

	for _, rule := range rules {
		for _, ruleTag := range rule.Tags {
			if ruleTag == tag {
				result = append(result, rule)
				break
			}
		}
	}

	return result
}
